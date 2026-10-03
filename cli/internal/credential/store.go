// Package credential stores and resolves the webmemo server URL and token.
package credential

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultServer = "https://memo.redshore.me"

	envToken  = "WEBMEMO_TOKEN"
	envServer = "WEBMEMO_SERVER"
)

type Credential struct {
	Server string `json:"server"`
	Token  string `json:"token"`
}

// Source tells where a resolved token came from.
type Source int

const (
	SourceNone Source = iota
	SourceEnv
	SourceFile
)

// Store persists a Credential as a JSON file.
type Store struct {
	Path string
}

// DefaultStore returns the store at $XDG_CONFIG_HOME/webmemo/credentials.json,
// falling back to ~/.config/webmemo/credentials.json.
func DefaultStore() (*Store, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("find home directory: %w", err)
		}
		base = filepath.Join(home, ".config")
	}
	return &Store{Path: filepath.Join(base, "webmemo", "credentials.json")}, nil
}

// Load reads the stored credential. A missing file yields an empty Credential.
func (s *Store) Load() (Credential, error) {
	data, err := os.ReadFile(s.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return Credential{}, nil
	}
	if err != nil {
		return Credential{}, fmt.Errorf("read credentials file %s: %w", s.Path, err)
	}
	var c Credential
	if err := json.Unmarshal(data, &c); err != nil {
		return Credential{}, fmt.Errorf("parse credentials file %s: %w", s.Path, err)
	}
	return c, nil
}

// Save writes the credential atomically (temp file + rename) with mode 0600,
// creating the parent directory with mode 0700.
func (s *Store) Save(c Credential) error {
	dir := filepath.Dir(s.Path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create credentials directory %s: %w", dir, err)
	}
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("encode credentials: %w", err)
	}

	tmp, err := os.CreateTemp(dir, ".credentials-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp credentials file in %s: %w", dir, err)
	}
	tmpPath := tmp.Name()
	cleanup := func() { _ = os.Remove(tmpPath) }

	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("chmod temp credentials file: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("write temp credentials file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return fmt.Errorf("close temp credentials file: %w", err)
	}
	if err := os.Rename(tmpPath, s.Path); err != nil {
		cleanup()
		return fmt.Errorf("replace credentials file %s: %w", s.Path, err)
	}
	return nil
}

// Delete removes the credentials file. A missing file is not an error.
func (s *Store) Delete() error {
	if err := os.Remove(s.Path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove credentials file %s: %w", s.Path, err)
	}
	return nil
}

// SameServer reports whether two server URLs are the same apart from
// trailing slashes. An empty URL means DefaultServer.
func SameServer(a, b string) bool {
	norm := func(u string) string {
		if u == "" {
			u = DefaultServer
		}
		return strings.TrimRight(u, "/")
	}
	return norm(a) == norm(b)
}

// Resolve combines environment and file credentials. The server is taken from
// WEBMEMO_SERVER, then the file, then DefaultServer. The token is taken from
// WEBMEMO_TOKEN (SourceEnv), then the file (SourceFile), else SourceNone.
//
// The file token is not used when WEBMEMO_SERVER names a different server than
// the file does, so a token is never sent to a server it was not issued by.
// An unreadable file is an error unless WEBMEMO_TOKEN is set, in which case
// the file is ignored.
func Resolve(s *Store, getenv func(string) string) (Credential, Source, error) {
	envTok := getenv(envToken)
	file, err := s.Load()
	if err != nil {
		if envTok == "" {
			return Credential{}, SourceNone, err
		}
		file = Credential{}
	}

	server := DefaultServer
	if file.Server != "" {
		server = file.Server
	}
	envSrv := getenv(envServer)
	if envSrv != "" {
		server = envSrv
	}

	switch {
	case envTok != "":
		return Credential{Server: server, Token: envTok}, SourceEnv, nil
	case file.Token != "" && (envSrv == "" || SameServer(envSrv, file.Server)):
		return Credential{Server: server, Token: file.Token}, SourceFile, nil
	default:
		return Credential{Server: server}, SourceNone, nil
	}
}
