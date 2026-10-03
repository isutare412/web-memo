package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/isutare412/web-memo/cli/internal/credential"
)

func noEnv(string) string { return "" }

func makeJWT(exp int64) string {
	enc := base64.RawURLEncoding
	return enc.EncodeToString([]byte(`{"alg":"none"}`)) + "." +
		enc.EncodeToString([]byte(fmt.Sprintf(`{"exp":%d}`, exp))) + ".sig"
}

func newStore(t *testing.T) *credential.Store {
	t.Helper()
	return &credential.Store{Path: filepath.Join(t.TempDir(), "webmemo", "credentials.json")}
}

// fakeAPI serves /api/v1/users/me, accepting only wantToken.
func fakeAPI(t *testing.T, wantToken string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/users/me" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("Authorization") != "Bearer "+wantToken {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"msg":"unauthorized"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"11111111-1111-1111-1111-111111111111","email":"jane@example.com","userName":"jane","userType":"user","issuedAt":"2026-01-01T00:00:00Z","expireAt":"2030-01-01T00:00:00Z"}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestLoginWithToken(t *testing.T) {
	tok := makeJWT(1893456000)
	srv := fakeAPI(t, tok)
	store := newStore(t)
	var out strings.Builder

	err := runLogin(context.Background(), []string{"--server", srv.URL, "--token", tok}, noEnv, store, &out)
	if err != nil {
		t.Fatalf("runLogin: %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Server != srv.URL || got.Token != tok {
		t.Fatalf("stored credential = %+v", got)
	}
	if strings.Contains(out.String(), tok) {
		t.Fatalf("output leaks the token: %q", out.String())
	}
	if !strings.Contains(out.String(), "jane@example.com") {
		t.Fatalf("output should name the user: %q", out.String())
	}
}

func TestLoginWithTokenRejected(t *testing.T) {
	srv := fakeAPI(t, "good")
	store := newStore(t)

	err := runLogin(context.Background(), []string{"--server", srv.URL, "--token", "bad-secret"}, noEnv, store, &strings.Builder{})
	if err == nil {
		t.Fatal("runLogin succeeded with a rejected token")
	}
	if strings.Contains(err.Error(), "bad-secret") {
		t.Fatalf("error leaks the token: %v", err)
	}
	if _, statErr := os.Stat(store.Path); !os.IsNotExist(statErr) {
		t.Fatalf("credentials file must not exist, stat err = %v", statErr)
	}
}

func TestLoginWithCorruptFile(t *testing.T) {
	tok := makeJWT(1893456000)
	srv := fakeAPI(t, tok)
	store := newStore(t)
	if err := os.MkdirAll(filepath.Dir(store.Path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(store.Path, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	env := func(k string) string {
		if k == "WEBMEMO_SERVER" {
			return srv.URL
		}
		return ""
	}

	if err := runLogin(context.Background(), []string{"--token", tok}, env, store, &out); err != nil {
		t.Fatalf("runLogin: %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("file should be rewritten: %v", err)
	}
	if got.Server != srv.URL || got.Token != tok {
		t.Fatalf("stored credential = %+v", got)
	}
	if !strings.Contains(out.String(), "unreadable") {
		t.Fatalf("output should mention the unreadable file: %q", out.String())
	}
}

func TestLoginNotesEnvToken(t *testing.T) {
	tok := makeJWT(1893456000)
	srv := fakeAPI(t, tok)
	var out strings.Builder
	env := func(k string) string {
		if k == "WEBMEMO_TOKEN" {
			return "other"
		}
		return ""
	}

	err := runLogin(context.Background(), []string{"--server", srv.URL, "--token", tok}, env, newStore(t), &out)
	if err != nil {
		t.Fatalf("runLogin: %v", err)
	}
	if !strings.Contains(out.String(), "WEBMEMO_TOKEN is set") {
		t.Fatalf("output should note WEBMEMO_TOKEN: %q", out.String())
	}
	if strings.Contains(out.String(), "other") {
		t.Fatalf("output leaks the env token: %q", out.String())
	}
}

func TestLoginBrowser(t *testing.T) {
	tok := makeJWT(1893456000)
	srv := fakeAPI(t, tok)
	store := newStore(t)

	prev := browserOpener
	t.Cleanup(func() { browserOpener = prev })
	browserOpener = func(raw string) error {
		u, err := url.Parse(raw)
		if err != nil {
			return err
		}
		q := u.Query()
		go func() {
			resp, err := http.Get(q.Get("cliCallback") + "?token=" + tok + "&state=" + q.Get("cliState"))
			if err == nil {
				_ = resp.Body.Close()
			}
		}()
		return nil
	}

	var out strings.Builder
	if err := runLogin(context.Background(), []string{"--server", srv.URL}, noEnv, store, &out); err != nil {
		t.Fatalf("runLogin: %v", err)
	}
	got, _ := store.Load()
	if got.Server != srv.URL || got.Token != tok {
		t.Fatalf("stored credential = %+v", got)
	}
	if strings.Contains(out.String(), tok) {
		t.Fatalf("output leaks the token: %q", out.String())
	}
}

func TestLoginBrowserRejectedTokenNotSaved(t *testing.T) {
	srv := fakeAPI(t, "good")
	store := newStore(t)

	prev := browserOpener
	t.Cleanup(func() { browserOpener = prev })
	browserOpener = func(raw string) error {
		u, _ := url.Parse(raw)
		q := u.Query()
		go func() {
			resp, err := http.Get(q.Get("cliCallback") + "?token=bad&state=" + q.Get("cliState"))
			if err == nil {
				_ = resp.Body.Close()
			}
		}()
		return nil
	}
	if err := runLogin(context.Background(), []string{"--server", srv.URL}, noEnv, store, &strings.Builder{}); err == nil {
		t.Fatal("runLogin succeeded with a rejected token")
	}
	if _, statErr := os.Stat(store.Path); !os.IsNotExist(statErr) {
		t.Fatalf("credentials file must not exist, stat err = %v", statErr)
	}
}

func TestLoginNoBrowserPrintsURL(t *testing.T) {
	srv := fakeAPI(t, "x")
	store := newStore(t)

	prev := browserOpener
	t.Cleanup(func() { browserOpener = prev })
	browserOpener = func(string) error {
		t.Error("browser must not be opened with --no-browser")
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	var out strings.Builder
	_ = runLogin(ctx, []string{"--server", srv.URL, "--no-browser"}, noEnv, store, &out)
	if !strings.Contains(out.String(), "/api/v1/google/sign-in") {
		t.Fatalf("login URL not printed: %q", out.String())
	}
}

func TestTokenCommandPrintsToken(t *testing.T) {
	store := newStore(t)
	if err := store.Save(credential.Credential{Server: "https://s.example", Token: "stored-tok"}); err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := runToken(nil, noEnv, store, &out); err != nil {
		t.Fatalf("runToken: %v", err)
	}
	if out.String() != "stored-tok\n" {
		t.Fatalf("output = %q", out.String())
	}

	env := func(k string) string {
		if k == "WEBMEMO_TOKEN" {
			return "env-tok"
		}
		return ""
	}
	out.Reset()
	if err := runToken(nil, env, store, &out); err != nil {
		t.Fatal(err)
	}
	if out.String() != "env-tok\n" {
		t.Fatalf("env output = %q", out.String())
	}
}

func TestTokenCommandNotLoggedIn(t *testing.T) {
	if err := runToken(nil, noEnv, newStore(t), &strings.Builder{}); err == nil {
		t.Fatal("runToken succeeded without credentials")
	}
}

func TestLogoutDeletes(t *testing.T) {
	store := newStore(t)
	if err := store.Save(credential.Credential{Server: "s", Token: "t"}); err != nil {
		t.Fatal(err)
	}
	if err := runLogout(nil, noEnv, store, &strings.Builder{}); err != nil {
		t.Fatalf("runLogout: %v", err)
	}
	if _, err := os.Stat(store.Path); !os.IsNotExist(err) {
		t.Fatalf("credentials file still exists, stat err = %v", err)
	}
	// Logging out twice is fine.
	if err := runLogout(nil, noEnv, store, &strings.Builder{}); err != nil {
		t.Fatalf("second runLogout: %v", err)
	}
}

func TestWhoamiPrintsUserAndExpiry(t *testing.T) {
	tok := makeJWT(1893456000)
	srv := fakeAPI(t, tok)
	store := newStore(t)
	if err := store.Save(credential.Credential{Server: srv.URL, Token: tok}); err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := runWhoami(context.Background(), nil, noEnv, store, &out); err != nil {
		t.Fatalf("runWhoami: %v", err)
	}
	s := out.String()
	wantExp := time.Unix(1893456000, 0).Format(time.RFC3339)
	for _, want := range []string{"jane", "jane@example.com", wantExp} {
		if !strings.Contains(s, want) {
			t.Errorf("output %q missing %q", s, want)
		}
	}
	if strings.Contains(s, tok) {
		t.Errorf("output leaks the token: %q", s)
	}
}

func TestWhoamiUnauthorized(t *testing.T) {
	srv := fakeAPI(t, "good")
	store := newStore(t)
	if err := store.Save(credential.Credential{Server: srv.URL, Token: "stale"}); err != nil {
		t.Fatal(err)
	}
	if err := runWhoami(context.Background(), nil, noEnv, store, &strings.Builder{}); err == nil {
		t.Fatal("runWhoami succeeded with a rejected token")
	}
}

func TestRunUnknownSubcommand(t *testing.T) {
	var out strings.Builder
	if err := run(context.Background(), []string{"bogus"}, noEnv, newStore(t), &out); err == nil {
		t.Fatal("run succeeded with an unknown subcommand")
	}
	if err := run(context.Background(), nil, noEnv, newStore(t), &out); err == nil {
		t.Fatal("run succeeded with no subcommand")
	}
}
