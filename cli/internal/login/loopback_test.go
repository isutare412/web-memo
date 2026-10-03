package login_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/isutare412/web-memo/cli/internal/login"
)

func get(t *testing.T, rawURL string) (int, string) {
	t.Helper()
	resp, err := http.Get(rawURL)
	if err != nil {
		t.Errorf("GET %s: %v", rawURL, err)
		return 0, ""
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

func parseLoginURL(t *testing.T, raw string) (callback, state string) {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Errorf("parse login url: %v", err)
		return "", ""
	}
	if u.Path != "/api/v1/google/sign-in" {
		t.Errorf("login url path = %q, want /api/v1/google/sign-in", u.Path)
	}
	q := u.Query()
	return q.Get("cliCallback"), q.Get("cliState")
}

func TestLoopbackSuccess(t *testing.T) {
	var wg sync.WaitGroup
	var gotStatus int
	var gotBody string
	lb := &login.Loopback{
		Server:  "https://memo.example.test",
		Timeout: 5 * time.Second,
		OpenBrowser: func(raw string) error {
			cb, state := parseLoginURL(t, raw)
			if len(state) != 43 {
				t.Errorf("state length = %d, want 43", len(state))
			}
			cbURL, err := url.Parse(cb)
			if err != nil || cbURL.Hostname() != "127.0.0.1" || cbURL.Path != "/callback" {
				t.Errorf("unexpected cliCallback %q (err %v)", cb, err)
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				gotStatus, gotBody = get(t, cb+"?token=tok&state="+url.QueryEscape(state))
			}()
			return nil
		},
	}
	var out strings.Builder
	tok, err := lb.Run(context.Background(), &out)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	wg.Wait()
	if tok != "tok" {
		t.Fatalf("token = %q, want tok", tok)
	}
	if gotStatus != http.StatusOK {
		t.Fatalf("callback status = %d, want 200", gotStatus)
	}
	if !strings.Contains(gotBody, "webmemo 로그인 완료") {
		t.Fatalf("callback body = %q", gotBody)
	}
	if !strings.Contains(out.String(), "/api/v1/google/sign-in") {
		t.Fatalf("login URL not printed: %q", out.String())
	}
	if strings.Contains(out.String(), "tok\n") {
		t.Fatalf("output must not contain the token: %q", out.String())
	}
	if !strings.Contains(out.String(), "Ctrl-C") {
		t.Fatalf("output should tell how to cancel: %q", out.String())
	}
}

func TestLoopbackIgnoresStray(t *testing.T) {
	goodSent := make(chan struct{})
	done := make(chan struct{})
	lb := &login.Loopback{
		Server:  "https://memo.example.test",
		Timeout: 5 * time.Second,
		OpenBrowser: func(raw string) error {
			cb, state := parseLoginURL(t, raw)
			base := strings.TrimSuffix(cb, "/callback")
			go func() {
				defer close(done)
				if s, _ := get(t, base+"/favicon.ico"); s != http.StatusNotFound {
					t.Errorf("favicon status = %d, want 404", s)
				}
				if s, _ := get(t, cb+"?token=evil&state=wrong"); s != http.StatusBadRequest {
					t.Errorf("state mismatch status = %d, want 400", s)
				}
				if s, _ := get(t, cb+"?state="+state); s != http.StatusBadRequest {
					t.Errorf("missing token status = %d, want 400", s)
				}
				close(goodSent)
				if s, _ := get(t, cb+"?token=good&state="+state); s != http.StatusOK {
					t.Errorf("good status = %d, want 200", s)
				}
			}()
			return nil
		},
	}
	tok, err := lb.Run(context.Background(), io.Discard)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	select {
	case <-goodSent:
	default:
		t.Fatal("Run returned before the valid callback was sent")
	}
	<-done
	if tok != "good" {
		t.Fatalf("token = %q, want good", tok)
	}
}

func TestLoopbackTimeout(t *testing.T) {
	lb := &login.Loopback{Server: "https://memo.example.test", Timeout: 50 * time.Millisecond}
	var out strings.Builder
	if _, err := lb.Run(context.Background(), &out); err == nil {
		t.Fatal("Run succeeded, want timeout error")
	}
	if !strings.Contains(out.String(), "cliCallback=") {
		t.Fatalf("with nil OpenBrowser the URL must still be printed: %q", out.String())
	}
}

func TestLoopbackContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	lb := &login.Loopback{Server: "https://memo.example.test", Timeout: 5 * time.Second}
	if _, err := lb.Run(ctx, io.Discard); err == nil {
		t.Fatal("Run succeeded, want context error")
	}
}
