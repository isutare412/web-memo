// Package login implements the browser-based login flow: the CLI listens on a
// loopback port, sends the user to the server's Google sign-in with that port
// as the callback, and receives the app token on redirect.
package login

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout = 5 * time.Minute
	signInPath     = "/api/v1/google/sign-in"
	callbackPath   = "/callback"
	successPage    = "<!doctype html><meta charset=\"utf-8\"><title>webmemo</title><p>webmemo 로그인 완료. 이 창을 닫아도 됩니다.</p>\n"
)

// Loopback runs the loopback login flow against Server.
type Loopback struct {
	// Server is the base URL of the webmemo server.
	Server string
	// OpenBrowser opens the login URL. If nil, the URL is only printed.
	OpenBrowser func(url string) error
	// Timeout bounds the wait for the callback. Zero means five minutes.
	Timeout time.Duration
}

// Run starts the callback listener, prints the login URL to out, opens the
// browser, and waits for the token. It returns an error on timeout or when ctx
// is canceled. Requests that are not a valid callback are rejected and ignored.
func (l *Loopback) Run(ctx context.Context, out io.Writer) (string, error) {
	timeout := l.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	state, err := newState()
	if err != nil {
		return "", err
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("listen on loopback: %w", err)
	}

	tokens := make(chan string, 1)
	srv := &http.Server{
		Handler:           callbackHandler(state, tokens),
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() { _ = srv.Serve(ln) }()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			_ = srv.Close()
		}
	}()

	loginURL, err := l.loginURL(ln.Addr().(*net.TCPAddr).Port, state)
	if err != nil {
		return "", err
	}
	_, _ = fmt.Fprintf(out, "Open this URL in your browser to log in:\n\n  %s\n\n", loginURL)
	if l.OpenBrowser != nil {
		if err := l.OpenBrowser(loginURL); err != nil {
			_, _ = fmt.Fprintf(out, "Could not open a browser automatically (%v); open the URL above manually.\n", err)
		}
	}
	_, _ = fmt.Fprintln(out, "Waiting for login to complete... (if the browser shows an error or you cancel the consent, press Ctrl-C)")

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case tok := <-tokens:
		return tok, nil
	case <-timer.C:
		return "", fmt.Errorf("login timed out after %s", timeout)
	case <-ctx.Done():
		return "", fmt.Errorf("login canceled: %w", ctx.Err())
	}
}

func (l *Loopback) loginURL(port int, state string) (string, error) {
	base, err := url.Parse(strings.TrimRight(l.Server, "/") + signInPath)
	if err != nil {
		return "", fmt.Errorf("parse server url: %w", err)
	}
	q := url.Values{}
	q.Set("cliCallback", fmt.Sprintf("http://127.0.0.1:%d%s", port, callbackPath))
	q.Set("cliState", state)
	base.RawQuery = q.Encode()
	return base.String(), nil
}

func callbackHandler(state string, tokens chan<- string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(callbackPath, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if subtle.ConstantTimeCompare([]byte(q.Get("state")), []byte(state)) != 1 {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}
		tok := q.Get("token")
		if tok == "" {
			http.Error(w, "missing token", http.StatusBadRequest)
			return
		}
		select {
		case tokens <- tok:
		default:
			// A token was already received; ignore replays.
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = io.WriteString(w, successPage)
	})
	return mux
}

func newState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate login state: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
