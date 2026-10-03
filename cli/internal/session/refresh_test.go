package session_test

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/isutare412/web-memo/cli/internal/credential"
	"github.com/isutare412/web-memo/cli/internal/session"
)

func TestForceRefreshFarFromExpiry(t *testing.T) {
	f := newFakeAPI(t, "newTok", http.StatusOK)
	old := jwtExpiringIn(29 * 24 * time.Hour)
	store := newFileStore(t)
	cred := credential.Credential{Server: f.srv.URL, Token: old}
	if err := store.Save(cred); err != nil {
		t.Fatal(err)
	}
	s := session.New(cred, credential.SourceFile, store, nil)

	tok, err := s.Refresh(context.Background())
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if tok != "newTok" || s.Token() != "newTok" {
		t.Fatalf("token = %q / %q, want newTok", tok, s.Token())
	}
	if f.hits() != 1 || f.refreshAuth != "Bearer "+old {
		t.Fatalf("refresh hits = %d auth = %q", f.hits(), f.refreshAuth)
	}
	saved, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if saved.Token != "newTok" || saved.Server != f.srv.URL {
		t.Fatalf("saved = %+v", saved)
	}
}

func TestForceRefreshEnvSourceNotSaved(t *testing.T) {
	f := newFakeAPI(t, "newTok", http.StatusOK)
	store := newFileStore(t)
	s := session.New(credential.Credential{Server: f.srv.URL, Token: "envTok"}, credential.SourceEnv, store, nil)

	tok, err := s.Refresh(context.Background())
	if err != nil {
		t.Fatalf("Refresh: %v", err)
	}
	if tok != "newTok" {
		t.Fatalf("token = %q", tok)
	}
	if _, err := os.Stat(store.Path); !os.IsNotExist(err) {
		t.Fatalf("credentials file written for env source (stat err = %v)", err)
	}
}

func TestForceRefreshFailures(t *testing.T) {
	t.Run("server error", func(t *testing.T) {
		f := newFakeAPI(t, "", http.StatusInternalServerError)
		s := session.New(credential.Credential{Server: f.srv.URL, Token: "tok"}, credential.SourceEnv, nil, nil)
		if _, err := s.Refresh(context.Background()); err == nil {
			t.Fatal("Refresh succeeded on 500")
		}
		if s.Token() != "tok" {
			t.Fatalf("token changed to %q", s.Token())
		}
	})
	t.Run("unauthorized", func(t *testing.T) {
		f := newFakeAPI(t, "", http.StatusUnauthorized)
		s := session.New(credential.Credential{Server: f.srv.URL, Token: "tok"}, credential.SourceEnv, nil, nil)
		if _, err := s.Refresh(context.Background()); !errors.Is(err, session.ErrUnauthorized) {
			t.Fatalf("err = %v, want ErrUnauthorized", err)
		}
	})
	t.Run("no token", func(t *testing.T) {
		f := newFakeAPI(t, "newTok", http.StatusOK)
		s := session.New(credential.Credential{Server: f.srv.URL}, credential.SourceNone, nil, nil)
		if _, err := s.Refresh(context.Background()); !errors.Is(err, session.ErrUnauthorized) {
			t.Fatalf("err = %v, want ErrUnauthorized", err)
		}
		if f.hits() != 0 {
			t.Fatalf("refresh endpoint called without a token")
		}
	})
}
