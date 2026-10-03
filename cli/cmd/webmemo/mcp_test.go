package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/isutare412/web-memo/cli/internal/credential"
)

// connectMCP builds the server for args and connects an in-memory client.
func connectMCP(t *testing.T, args []string, env func(string) string) *mcp.ClientSession {
	t.Helper()
	return connectMCPStore(t, args, env, newStore(t))
}

func connectMCPStore(t *testing.T, args []string, env func(string) string, store *credential.Store) *mcp.ClientSession {
	t.Helper()
	srv, err := newMCPServer(args, env, store)
	if err != nil {
		t.Fatalf("newMCPServer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := srv.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("connect server: %v", err)
	}
	t.Cleanup(func() { _ = serverSession.Close() })
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0"}, nil)
	cs, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("connect client: %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

func envOf(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestMcpCommandBuildsServer(t *testing.T) {
	var mu sync.Mutex
	var gotAuth string
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		gotAuth = r.Header.Get("Authorization")
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"memos":[],"paging":{"page":1,"pageSize":10,"lastPage":1,"totalMemoCount":0}}`))
	}))
	t.Cleanup(api.Close)

	cs := connectMCP(t, nil, envOf(map[string]string{"WEBMEMO_TOKEN": "env-token", "WEBMEMO_SERVER": api.URL}))

	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "search_memos", Arguments: map[string]any{}})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool error: %+v", res.Content)
	}
	mu.Lock()
	defer mu.Unlock()
	if gotAuth != "Bearer env-token" {
		t.Errorf("Authorization = %q, want Bearer env-token", gotAuth)
	}

	tools, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 14 {
		t.Errorf("tools = %d, want 14", len(tools.Tools))
	}
}

func TestMcpCommandReadOnly(t *testing.T) {
	cs := connectMCP(t, []string{"--read-only"}, noEnv)
	tools, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 4 {
		t.Errorf("tools = %d, want 4", len(tools.Tools))
	}
}

func TestMcpCommandStartsWithoutToken(t *testing.T) {
	cs := connectMCP(t, nil, noEnv)
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "search_memos", Arguments: map[string]any{}})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if !res.IsError {
		t.Fatal("IsError = false, want an unauthorized tool error")
	}
	text := res.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(text, "webmemo login") {
		t.Errorf("text = %q, want login guidance", text)
	}
}

func TestMcpCommandRejectsArgs(t *testing.T) {
	if _, err := newMCPServer([]string{"extra"}, noEnv, newStore(t)); err == nil {
		t.Error("want error for positional argument")
	}
}

func TestMcpPicksUpLoginWhileRunning(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fresh" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"memos":[],"paging":{"page":1,"pageSize":10,"lastPage":1,"totalMemoCount":0}}`))
	}))
	t.Cleanup(api.Close)
	store := newStore(t)
	cs := connectMCPStore(t, nil, envOf(map[string]string{"WEBMEMO_SERVER": api.URL}), store)
	call := func() *mcp.CallToolResult {
		res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "search_memos", Arguments: map[string]any{}})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		return res
	}

	if !call().IsError {
		t.Fatal("first call should fail: not logged in")
	}
	if err := store.Save(credential.Credential{Server: api.URL, Token: "fresh"}); err != nil {
		t.Fatal(err)
	}
	if res := call(); res.IsError {
		t.Fatalf("call after login failed: %+v", res.Content)
	}
}
