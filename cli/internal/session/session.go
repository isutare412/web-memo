// Package session provides an authenticated API client that sends the Bearer
// token and refreshes it before it expires.
package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
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
	next   http.RoundTripper
	src    credential.Source
	store  *credential.Store
	now    func() time.Time
	apiURL string

	mu         sync.Mutex
	cred       credential.Credential
	retryAfter time.Time
}

func (t *authTransport) token() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.cred.Token
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var tok string
	if req.Method == http.MethodPost && req.URL.Path == refreshPath {
		tok = t.token()
	} else {
		tok = t.currentToken(req)
	}
	req = req.Clone(req.Context())
	if tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	return t.next.RoundTrip(req)
}

// currentToken returns the token to use, refreshing it first when it is close
// to expiry. The lock is held during the refresh so concurrent requests wait
// for it instead of refreshing again.
func (t *authTransport) currentToken(req *http.Request) string {
	t.mu.Lock()
	defer t.mu.Unlock()

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
