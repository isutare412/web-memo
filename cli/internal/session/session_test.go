package session_test

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/isutare412/web-memo/cli/internal/credential"
	"github.com/isutare412/web-memo/cli/internal/session"
)

func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	os.Exit(m.Run())
}

func jwtExpiringIn(d time.Duration) string {
	enc := base64.RawURLEncoding
	payload := fmt.Sprintf(`{"exp":%d}`, time.Now().Add(d).Unix())
	return enc.EncodeToString([]byte(`{"alg":"none"}`)) + "." + enc.EncodeToString([]byte(payload)) + ".sig"
}

// fakeAPI records the Authorization header of each request and serves the
// refresh endpoint.
type fakeAPI struct {
	srv *httptest.Server

	mu          sync.Mutex
	refreshHits int
	memoAuths   []string
	refreshAuth string
	newToken    string // returned in Set-Cookie when non-empty and refreshStatus is 200
	refreshCode int
}

func newFakeAPI(t *testing.T, newToken string, refreshCode int) *fakeAPI {
	t.Helper()
	f := &fakeAPI{newToken: newToken, refreshCode: refreshCode}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/users/refresh-token", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.refreshHits++
		f.refreshAuth = r.Header.Get("Authorization")
		f.mu.Unlock()
		if f.refreshCode != http.StatusOK {
			w.WriteHeader(f.refreshCode)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "wmToken", Value: f.newToken, Path: "/"})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})
	mux.HandleFunc("GET /api/v1/users/me", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.memoAuths = append(f.memoAuths, r.Header.Get("Authorization"))
		f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})
	f.srv = httptest.NewServer(mux)
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeAPI) auths() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.memoAuths...)
}

func (f *fakeAPI) hits() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.refreshHits
}

func callAPI(t *testing.T, s *session.Session) {
	t.Helper()
	if _, err := s.API().GetCurrentUserWithResponse(context.Background()); err != nil {
		t.Fatalf("api call: %v", err)
	}
}

func TestBearerHeader(t *testing.T) {
	f := newFakeAPI(t, "", http.StatusOK)
	tok := jwtExpiringIn(30 * 24 * time.Hour)
	s := session.New(credential.Credential{Server: f.srv.URL, Token: tok}, credential.SourceEnv, nil, nil)

	callAPI(t, s)

	got := f.auths()
	if len(got) != 1 || got[0] != "Bearer "+tok {
		t.Fatalf("auth headers = %v, want [Bearer <token>]", got)
	}
	if s.Token() != tok {
		t.Fatalf("Token() changed unexpectedly")
	}
}

func TestRefreshNearExpiry(t *testing.T) {
	t.Run("file source persists new token", func(t *testing.T) {
		f := newFakeAPI(t, "newTok", http.StatusOK)
		old := jwtExpiringIn(3 * 24 * time.Hour)
		store := &credential.Store{Path: filepath.Join(t.TempDir(), "webmemo", "credentials.json")}
		cred := credential.Credential{Server: f.srv.URL, Token: old}
		if err := store.Save(cred); err != nil {
			t.Fatal(err)
		}
		s := session.New(cred, credential.SourceFile, store, nil)

		callAPI(t, s)
		callAPI(t, s)

		if f.hits() != 1 {
			t.Fatalf("refresh hits = %d, want 1", f.hits())
		}
		if f.refreshAuth != "Bearer "+old {
			t.Fatalf("refresh Authorization = %q, want old bearer", f.refreshAuth)
		}
		for i, a := range f.auths() {
			if a != "Bearer newTok" {
				t.Fatalf("request %d auth = %q, want Bearer newTok", i, a)
			}
		}
		if s.Token() != "newTok" {
			t.Fatalf("Token() = %q, want newTok", s.Token())
		}
		saved, err := store.Load()
		if err != nil {
			t.Fatal(err)
		}
		if saved.Token != "newTok" || saved.Server != f.srv.URL {
			t.Fatalf("saved credential = %+v", saved)
		}
	})

	t.Run("env source does not write file", func(t *testing.T) {
		f := newFakeAPI(t, "newTok", http.StatusOK)
		path := filepath.Join(t.TempDir(), "webmemo", "credentials.json")
		store := &credential.Store{Path: path}
		s := session.New(credential.Credential{Server: f.srv.URL, Token: jwtExpiringIn(3 * 24 * time.Hour)},
			credential.SourceEnv, store, nil)

		callAPI(t, s)

		if got := f.auths(); len(got) != 1 || got[0] != "Bearer newTok" {
			t.Fatalf("auth headers = %v", got)
		}
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("credentials file should not exist, stat err = %v", err)
		}
	})

	t.Run("nil store with env source", func(t *testing.T) {
		f := newFakeAPI(t, "newTok", http.StatusOK)
		s := session.New(credential.Credential{Server: f.srv.URL, Token: jwtExpiringIn(time.Hour)},
			credential.SourceEnv, nil, nil)
		callAPI(t, s)
		if s.Token() != "newTok" {
			t.Fatalf("Token() = %q, want newTok", s.Token())
		}
	})
}

func TestNoRefreshFarExpiry(t *testing.T) {
	f := newFakeAPI(t, "newTok", http.StatusOK)
	s := session.New(credential.Credential{Server: f.srv.URL, Token: jwtExpiringIn(20 * 24 * time.Hour)},
		credential.SourceEnv, nil, nil)

	callAPI(t, s)

	if f.hits() != 0 {
		t.Fatalf("refresh hits = %d, want 0", f.hits())
	}
}

func TestNoRefreshUnparsable(t *testing.T) {
	f := newFakeAPI(t, "newTok", http.StatusOK)
	s := session.New(credential.Credential{Server: f.srv.URL, Token: "garbage"}, credential.SourceEnv, nil, nil)

	callAPI(t, s)

	if f.hits() != 0 {
		t.Fatalf("refresh hits = %d, want 0", f.hits())
	}
	if got := f.auths(); len(got) != 1 || got[0] != "Bearer garbage" {
		t.Fatalf("auth headers = %v", got)
	}
}

func TestRefreshFailureContinues(t *testing.T) {
	f := newFakeAPI(t, "", http.StatusInternalServerError)
	old := jwtExpiringIn(3 * 24 * time.Hour)
	s := session.New(credential.Credential{Server: f.srv.URL, Token: old}, credential.SourceEnv, nil, nil)

	callAPI(t, s)
	callAPI(t, s)

	for i, a := range f.auths() {
		if a != "Bearer "+old {
			t.Fatalf("request %d auth = %q, want old bearer", i, a)
		}
	}
	if len(f.auths()) != 2 {
		t.Fatalf("requests = %d, want 2", len(f.auths()))
	}
	if f.hits() != 1 {
		t.Fatalf("refresh hits = %d, want 1 (no retry storm after failure)", f.hits())
	}
}

func TestRefreshWithoutCookieKeepsToken(t *testing.T) {
	f := newFakeAPI(t, "", http.StatusOK) // 200 but empty wmToken value
	old := jwtExpiringIn(3 * 24 * time.Hour)
	s := session.New(credential.Credential{Server: f.srv.URL, Token: old}, credential.SourceEnv, nil, nil)

	callAPI(t, s)

	if s.Token() != old {
		t.Fatalf("token replaced by empty cookie")
	}
}

func TestConcurrentRequestsRefreshOnce(t *testing.T) {
	f := newFakeAPI(t, "newTok", http.StatusOK)
	s := session.New(credential.Credential{Server: f.srv.URL, Token: jwtExpiringIn(time.Hour)},
		credential.SourceEnv, nil, nil)

	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() { _, _ = s.API().GetCurrentUserWithResponse(context.Background()) })
	}
	wg.Wait()

	if f.hits() != 1 {
		t.Fatalf("refresh hits = %d, want 1", f.hits())
	}
}

func TestHTTPClientIsUnauthenticated(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
	}))
	t.Cleanup(srv.Close)
	s := session.New(credential.Credential{Server: srv.URL, Token: "tok"}, credential.SourceEnv, nil, nil)

	resp, err := s.HTTPClient().Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if gotAuth != "" {
		t.Fatalf("Authorization = %q, want empty", gotAuth)
	}
}

func TestCheckStatus(t *testing.T) {
	if err := session.CheckStatus(204, nil); err != nil {
		t.Fatalf("2xx: %v", err)
	}
	if err := session.CheckStatus(401, []byte(`{"msg":"x"}`)); !errors.Is(err, session.ErrUnauthorized) {
		t.Fatalf("401: err = %v, want ErrUnauthorized", err)
	}
	err := session.CheckStatus(404, []byte(`{"msg":"memo not found"}`))
	if err == nil || err.Error() != "api error (404): memo not found" {
		t.Fatalf("404: err = %v", err)
	}
	err = session.CheckStatus(502, []byte(`<html>bad gateway</html>`))
	if err == nil || !strings.Contains(err.Error(), "api error (502)") {
		t.Fatalf("non-json body: err = %v", err)
	}
}

func TestRefreshRequestItselfIsNotRefreshed(t *testing.T) {
	f := newFakeAPI(t, "newTok", http.StatusOK)
	old := jwtExpiringIn(time.Hour)
	s := session.New(credential.Credential{Server: f.srv.URL, Token: old}, credential.SourceEnv, nil, nil)

	if _, err := s.API().RefreshTokenWithResponse(context.Background()); err != nil {
		t.Fatal(err)
	}

	if f.hits() != 1 {
		t.Fatalf("refresh hits = %d, want 1", f.hits())
	}
	if f.refreshAuth != "Bearer "+old {
		t.Fatalf("refresh Authorization = %q, want old bearer", f.refreshAuth)
	}
}
