package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/isutare412/web-memo/cli/internal/credential"
)

// refreshAPI serves the refresh endpoint, issuing newToken for wantToken.
func refreshAPI(t *testing.T, wantToken, newToken string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/users/refresh-token" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+wantToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "wmToken", Value: newToken, Path: "/"})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestRefreshCommandFileLogin(t *testing.T) {
	old, renewed := makeJWT(1893456000), makeJWT(1893999999)
	srv := refreshAPI(t, old, renewed)
	store := newStore(t)
	if err := store.Save(credential.Credential{Server: srv.URL, Token: old}); err != nil {
		t.Fatal(err)
	}
	var out, errOut strings.Builder

	if err := runRefresh(context.Background(), nil, noEnv, store, &out, &errOut); err != nil {
		t.Fatalf("runRefresh: %v", err)
	}
	saved, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if saved.Token != renewed {
		t.Fatalf("saved token not renewed")
	}
	if !strings.Contains(out.String(), "token refreshed; expires 2030-") {
		t.Fatalf("output = %q", out.String())
	}
	if strings.Contains(out.String()+errOut.String(), renewed) || strings.Contains(out.String()+errOut.String(), old) {
		t.Fatal("token printed for a file login")
	}
}

func TestRefreshCommandEnvToken(t *testing.T) {
	old, renewed := makeJWT(1893456000), makeJWT(1893999999)
	srv := refreshAPI(t, old, renewed)
	store := newStore(t)
	env := func(k string) string {
		switch k {
		case "WEBMEMO_TOKEN":
			return old
		case "WEBMEMO_SERVER":
			return srv.URL
		}
		return ""
	}
	var out, errOut strings.Builder

	if err := runRefresh(context.Background(), nil, env, store, &out, &errOut); err != nil {
		t.Fatalf("runRefresh: %v", err)
	}
	if out.String() != renewed+"\n" {
		t.Fatalf("stdout = %q, want only the new token", out.String())
	}
	if !strings.Contains(errOut.String(), "WEBMEMO_TOKEN") {
		t.Fatalf("stderr note missing: %q", errOut.String())
	}
	if c, _ := store.Load(); c.Token != "" {
		t.Fatal("env token refresh wrote the credentials file")
	}
}

func TestRefreshCommandErrors(t *testing.T) {
	t.Run("not logged in", func(t *testing.T) {
		if err := runRefresh(context.Background(), nil, noEnv, newStore(t), &strings.Builder{}, &strings.Builder{}); err == nil {
			t.Fatal("runRefresh succeeded without credentials")
		}
	})
	t.Run("rejected token", func(t *testing.T) {
		srv := refreshAPI(t, "other", "x")
		store := newStore(t)
		if err := store.Save(credential.Credential{Server: srv.URL, Token: "stale"}); err != nil {
			t.Fatal(err)
		}
		if err := runRefresh(context.Background(), nil, noEnv, store, &strings.Builder{}, &strings.Builder{}); err == nil {
			t.Fatal("runRefresh succeeded with a rejected token")
		}
		if c, _ := store.Load(); c.Token != "stale" {
			t.Fatal("stored token changed after a failed refresh")
		}
	})
}
