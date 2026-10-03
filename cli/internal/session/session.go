// Package session provides an authenticated API client that sends the Bearer
// token and refreshes it before it expires.
package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/isutare412/web-memo/cli/internal/apiclient/gen"
	"github.com/isutare412/web-memo/cli/internal/credential"
)

const (
	// refreshThreshold is how close to expiry a token must be to get refreshed.
	refreshThreshold = 7 * 24 * time.Hour
	// refreshRetryDelay keeps a failing refresh endpoint from being hit on
	// every request of a long-running process.
	refreshRetryDelay = time.Minute

	refreshPath = "/api/v1/users/refresh-token"
	cookieName  = "wmToken"
)

// ErrUnauthorized is returned when the server rejects the token or none is set.
// A running process picks up the new token after `webmemo login`.
var ErrUnauthorized = errors.New("webmemo: not logged in or token expired; run `webmemo login` (headless: `webmemo login --token <T>`)")

// Session holds the credential and builds API clients that authenticate with it.
type Session struct {
	base *http.Client
	api  *gen.ClientWithResponses
	tr   *authTransport
}

// New creates a Session. base is the unauthenticated HTTP client (nil means
// http.DefaultClient). A refreshed token is written to store only when src is
// SourceFile; store may be nil otherwise.
//
// Unless src is SourceEnv, a non-nil store is also watched for a newer token
// (for example after `webmemo login` in another terminal), so a long-running
// process picks up a new login without a restart. Only a token stored for the
// same server is adopted.
func New(cred credential.Credential, src credential.Source, store *credential.Store, base *http.Client) *Session {
	if base == nil {
		base = http.DefaultClient
	}
	rt := base.Transport
	if rt == nil {
		rt = http.DefaultTransport
	}
	tr := &authTransport{
		next:   rt,
		cred:   cred,
		src:    src,
		store:  store,
		now:    time.Now,
		apiURL: cred.Server,
	}
	if u, err := url.Parse(cred.Server); err == nil {
		tr.apiScheme, tr.apiHost = strings.ToLower(u.Scheme), strings.ToLower(u.Host)
	}
	if store != nil {
		tr.lastSig = fileSignature(store.Path)
	}
	authed := &http.Client{
		Transport:     tr,
		CheckRedirect: base.CheckRedirect,
		Jar:           base.Jar,
		Timeout:       base.Timeout,
	}
	api, err := gen.NewClientWithResponses(cred.Server, gen.WithHTTPClient(authed))
	if err != nil {
		// Unreachable: WithHTTPClient never fails and the server is not parsed.
		panic(fmt.Sprintf("session: build api client: %v", err))
	}
	return &Session{base: base, api: api, tr: tr}
}

// API returns the generated client that adds the Bearer header and refreshes
// the token near expiry.
func (s *Session) API() *gen.ClientWithResponses { return s.api }

// HTTPClient returns the unauthenticated client, e.g. for presigned uploads.
func (s *Session) HTTPClient() *http.Client { return s.base }

// Token returns the current token, which may have been refreshed.
func (s *Session) Token() string { return s.tr.token() }

// CheckStatus converts an HTTP status and body into an error: nil for 2xx,
// ErrUnauthorized for 401, otherwise an error carrying the API error message.
func CheckStatus(status int, body []byte) error {
	switch {
	case status >= 200 && status < 300:
		return nil
	case status == http.StatusUnauthorized:
		return ErrUnauthorized
	}
	var e gen.ErrorResponse
	if err := json.Unmarshal(body, &e); err != nil || e.Msg == "" {
		return fmt.Errorf("api error (%d): %s", status, http.StatusText(status))
	}
	return fmt.Errorf("api error (%d): %s", status, e.Msg)
}

type authTransport struct {
	next      http.RoundTripper
	src       credential.Source
	store     *credential.Store
	now       func() time.Time
	apiURL    string
	apiScheme string
	apiHost   string

	mu         sync.Mutex
	cred       credential.Credential
	retryAfter time.Time
	lastSig    fileSig
}

// fileSig is a cheap fingerprint of the credentials file.
type fileSig struct {
	exists  bool
	modTime time.Time
	size    int64
}

func fileSignature(path string) fileSig {
	info, err := os.Stat(path)
	if err != nil {
		return fileSig{}
	}
	return fileSig{exists: true, modTime: info.ModTime(), size: info.Size()}
}

func (t *authTransport) token() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.cred.Token
}

// isAPIHost reports whether u has the configured server's scheme and host.
// The token is only ever sent there, also after redirects.
func (t *authTransport) isAPIHost(u *url.URL) bool {
	return strings.EqualFold(u.Scheme, t.apiScheme) && strings.EqualFold(u.Host, t.apiHost)
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !t.isAPIHost(req.URL) {
		return t.next.RoundTrip(req)
	}

	isRefresh := req.Method == http.MethodPost && req.URL.Path == refreshPath
	var tok string
	if isRefresh {
		tok = t.token()
	} else {
		tok = t.currentToken(req)
	}
	resp, err := t.send(req, tok)
	if err != nil || resp.StatusCode != http.StatusUnauthorized || isRefresh {
		return resp, err
	}

	// The token may have been replaced by `webmemo login` since it was loaded.
	// Retry once, and only with a different token and a replayable body.
	newTok, changed := t.reloadAfterUnauthorized(tok)
	if !changed || (req.Body != nil && req.Body != http.NoBody && req.GetBody == nil) {
		return resp, nil
	}
	retry := req.Clone(req.Context())
	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return resp, nil
		}
		retry.Body = body
	}
	_ = resp.Body.Close()
	return t.send(retry, newTok)
}

func (t *authTransport) send(req *http.Request, tok string) (*http.Response, error) {
	req = req.Clone(req.Context())
	if tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	return t.next.RoundTrip(req)
}

// reloadFromFile adopts the token in the credentials file when it differs from
// the current one. Unless force is set it first compares a cheap file
// fingerprint and skips reading an unchanged file (a missing token is always
// worth a look). The caller must hold t.mu.
func (t *authTransport) reloadFromFile(force bool) {
	if t.src == credential.SourceEnv || t.store == nil {
		return
	}
	sig := fileSignature(t.store.Path)
	if !force && sig == t.lastSig && t.cred.Token != "" {
		return
	}
	t.lastSig = sig
	if !sig.exists {
		return
	}
	file, err := t.store.Load()
	if err != nil {
		slog.Warn("could not reload credentials", "error", err)
		return
	}
	if file.Token == "" || file.Token == t.cred.Token || !credential.SameServer(file.Server, t.cred.Server) {
		return
	}
	t.cred.Token = file.Token
	t.src = credential.SourceFile
	t.retryAfter = time.Time{}
}

// reloadAfterUnauthorized re-reads the credentials file after a 401 for sent
// and reports the current token and whether it differs from sent.
func (t *authTransport) reloadAfterUnauthorized(sent string) (string, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.reloadFromFile(true)
	return t.cred.Token, t.cred.Token != "" && t.cred.Token != sent
}

// currentToken returns the token to use, refreshing it first when it is close
// to expiry. The lock is held during the refresh so concurrent requests wait
// for it instead of refreshing again.
func (t *authTransport) currentToken(req *http.Request) string {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.reloadFromFile(false)

	exp, ok := credential.TokenExpiry(t.cred.Token)
	if !ok || exp.Sub(t.now()) >= refreshThreshold || t.now().Before(t.retryAfter) {
		return t.cred.Token
	}

	newTok, err := t.refresh(req)
	if err != nil {
		t.retryAfter = t.now().Add(refreshRetryDelay)
		slog.Warn("token refresh failed; continuing with current token", "error", err)
		return t.cred.Token
	}
	t.cred.Token = newTok
	if t.src == credential.SourceFile && t.store != nil {
		if err := t.store.Save(t.cred); err != nil {
			slog.Warn("could not save refreshed token", "error", err)
		}
	}
	return t.cred.Token
}

// refresh calls the refresh endpoint directly on the underlying transport, so
// refresh logic is never applied to the refresh request itself. The caller
// must hold t.mu.
func (t *authTransport) refresh(orig *http.Request) (string, error) {
	r, err := gen.NewRefreshTokenRequest(t.apiURL)
	if err != nil {
		return "", fmt.Errorf("build refresh request: %w", err)
	}
	r = r.WithContext(orig.Context())
	r.Header.Set("Authorization", "Bearer "+t.cred.Token)

	resp, err := t.next.RoundTrip(r)
	if err != nil {
		return "", fmt.Errorf("refresh request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("refresh request: status %d", resp.StatusCode)
	}
	for _, c := range resp.Cookies() {
		if c.Name == cookieName && c.Value != "" {
			return c.Value, nil
		}
	}
	return "", errors.New("refresh response has no token cookie")
}
