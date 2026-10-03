package credential

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return &Store{Path: filepath.Join(t.TempDir(), "webmemo", "credentials.json")}
}

func TestStoreSaveLoad(t *testing.T) {
	s := newTestStore(t)
	want := Credential{Server: "https://example.test", Token: "tok"}
	if err := s.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got != want {
		t.Fatalf("Load = %+v, want %+v", got, want)
	}
	info, err := os.Stat(s.Path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("file mode = %o, want 600", info.Mode().Perm())
	}
	dirInfo, err := os.Stat(filepath.Dir(s.Path))
	if err != nil {
		t.Fatal(err)
	}
	if dirInfo.Mode().Perm() != 0o700 {
		t.Fatalf("dir mode = %o, want 700", dirInfo.Mode().Perm())
	}
	entries, _ := os.ReadDir(filepath.Dir(s.Path))
	if len(entries) != 1 {
		t.Fatalf("expected only the credentials file, got %d entries", len(entries))
	}
}

func TestStoreSaveOverwrites(t *testing.T) {
	s := newTestStore(t)
	if err := s.Save(Credential{Server: "a", Token: "1"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Save(Credential{Server: "b", Token: "2"}); err != nil {
		t.Fatal(err)
	}
	got, err := s.Load()
	if err != nil || got != (Credential{Server: "b", Token: "2"}) {
		t.Fatalf("Load = %+v, %v", got, err)
	}
}

func TestStoreLoadMissing(t *testing.T) {
	got, err := newTestStore(t).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got != (Credential{}) {
		t.Fatalf("Load = %+v, want empty", got)
	}
}

func TestStoreLoadCorrupt(t *testing.T) {
	s := newTestStore(t)
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s.Path, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := s.Load()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), s.Path) {
		t.Fatalf("error %q does not contain path %q", err, s.Path)
	}
}

func TestStoreDelete(t *testing.T) {
	s := newTestStore(t)
	if err := s.Delete(); err != nil {
		t.Fatalf("Delete on missing: %v", err)
	}
	if err := s.Save(Credential{Server: "a", Token: "1"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(s.Path); !os.IsNotExist(err) {
		t.Fatalf("file still exists: %v", err)
	}
}

func TestDefaultStore(t *testing.T) {
	t.Run("xdg", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "/xdg")
		s, err := DefaultStore()
		if err != nil {
			t.Fatal(err)
		}
		if want := filepath.Join("/xdg", "webmemo", "credentials.json"); s.Path != want {
			t.Fatalf("Path = %q, want %q", s.Path, want)
		}
	})
	t.Run("home fallback", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("HOME", "/home/u")
		s, err := DefaultStore()
		if err != nil {
			t.Fatal(err)
		}
		if want := filepath.Join("/home/u", ".config", "webmemo", "credentials.json"); s.Path != want {
			t.Fatalf("Path = %q, want %q", s.Path, want)
		}
	})
}

func envFrom(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestResolvePrecedence(t *testing.T) {
	s := newTestStore(t)
	if err := s.Save(Credential{Server: "A", Token: "fileTok"}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		store      *Store
		env        map[string]string
		wantServer string
		wantToken  string
		wantSource Source
	}{
		{"file only", s, nil, "A", "fileTok", SourceFile},
		{"env token wins", s, map[string]string{"WEBMEMO_TOKEN": "envTok"}, "A", "envTok", SourceEnv},
		{"env server differing from file drops file token", s, map[string]string{"WEBMEMO_SERVER": "B"}, "B", "", SourceNone},
		{"env server equal to file keeps file token", s, map[string]string{"WEBMEMO_SERVER": "A/"}, "A/", "fileTok", SourceFile},
		{"env both", s, map[string]string{"WEBMEMO_TOKEN": "envTok", "WEBMEMO_SERVER": "B"}, "B", "envTok", SourceEnv},
		{"nothing", newTestStore(t), nil, DefaultServer, "", SourceNone},
		{"env server only", newTestStore(t), map[string]string{"WEBMEMO_SERVER": "B"}, "B", "", SourceNone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, src, err := Resolve(tt.store, envFrom(tt.env))
			if err != nil {
				t.Fatal(err)
			}
			if got.Server != tt.wantServer || got.Token != tt.wantToken || src != tt.wantSource {
				t.Fatalf("got (%q, %q, %v), want (%q, %q, %v)",
					got.Server, got.Token, src, tt.wantServer, tt.wantToken, tt.wantSource)
			}
		})
	}
}

func TestResolveDefaultServerFile(t *testing.T) {
	s := newTestStore(t)
	if err := s.Save(Credential{Token: "fileTok"}); err != nil { // no server stored: default
		t.Fatal(err)
	}
	got, src, err := Resolve(s, envFrom(map[string]string{"WEBMEMO_SERVER": DefaultServer}))
	if err != nil || src != SourceFile || got.Token != "fileTok" {
		t.Fatalf("got (%+v, %v, %v)", got, src, err)
	}
	_, src, err = Resolve(s, envFrom(map[string]string{"WEBMEMO_SERVER": "https://other.test"}))
	if err != nil || src != SourceNone {
		t.Fatalf("other server: src = %v, err = %v", src, err)
	}
}

func writeCorrupt(t *testing.T) *Store {
	t.Helper()
	s := newTestStore(t)
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s.Path, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestResolveCorruptFile(t *testing.T) {
	if _, _, err := Resolve(writeCorrupt(t), envFrom(nil)); err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveCorruptFileIgnoredWithEnvToken(t *testing.T) {
	got, src, err := Resolve(writeCorrupt(t), envFrom(map[string]string{"WEBMEMO_TOKEN": "envTok"}))
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if src != SourceEnv || got.Token != "envTok" || got.Server != DefaultServer {
		t.Fatalf("got (%+v, %v)", got, src)
	}
}
