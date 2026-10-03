package mcpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/isutare412/web-memo/cli/internal/credential"
	"github.com/isutare412/web-memo/cli/internal/session"
)

// newTestClient serves apiHandler over httptest, builds an MCP server backed by
// a session pointing at it, and returns a client connected over an in-memory
// transport.
func newTestClient(t *testing.T, apiHandler http.Handler, opts Options) *mcp.ClientSession {
	t.Helper()

	api := httptest.NewServer(apiHandler)
	t.Cleanup(api.Close)

	sess := session.New(
		credential.Credential{Server: api.URL, Token: "fake-token"},
		credential.SourceEnv,
		nil,
		&http.Client{Timeout: 10 * time.Second},
	)
	srv := New(sess, opts)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := srv.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("connect server: %v", err)
	}
	t.Cleanup(func() { _ = serverSession.Close() })

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("connect client: %v", err)
	}
	t.Cleanup(func() { _ = clientSession.Close() })
	return clientSession
}

func toolNames(t *testing.T, cs *mcp.ClientSession) []string {
	t.Helper()
	res, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	names := make([]string, 0, len(res.Tools))
	for _, tool := range res.Tools {
		names = append(names, tool.Name)
	}
	sort.Strings(names)
	return names
}

// resultText returns the text of the first content item.
func resultText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	if len(res.Content) == 0 {
		t.Fatalf("result has no content")
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("content[0] is %T, want *mcp.TextContent", res.Content[0])
	}
	return tc.Text
}

func callTool(t *testing.T, cs *mcp.ClientSession, name string, args any) *mcp.CallToolResult {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("call %s: %v", name, err)
	}
	return res
}

var readTools = []string{"get_memo", "list_memo_members", "search_memos", "search_tags"}

func TestToolListFull(t *testing.T) {
	cs := newTestClient(t, http.NotFoundHandler(), Options{})
	got := strings.Join(toolNames(t, cs), ",")
	for _, name := range readTools {
		if !strings.Contains(got, name) {
			t.Errorf("tool %q missing from %s", name, got)
		}
	}
}

func TestToolListReadOnly(t *testing.T) {
	cs := newTestClient(t, http.NotFoundHandler(), Options{ReadOnly: true})
	got := toolNames(t, cs)
	if strings.Join(got, ",") != strings.Join(readTools, ",") {
		t.Errorf("tools = %v, want %v", got, readTools)
	}
}

func TestReadToolsAreReadOnlyHinted(t *testing.T) {
	cs := newTestClient(t, http.NotFoundHandler(), Options{})
	res, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range res.Tools {
		if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
			t.Errorf("tool %s: ReadOnlyHint not set", tool.Name)
		}
		if tool.Description == "" {
			t.Errorf("tool %s: empty description", tool.Name)
		}
	}
}

func TestUnauthorized(t *testing.T) {
	cs := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}), Options{})

	res := callTool(t, cs, "search_memos", map[string]any{})
	if !res.IsError {
		t.Fatalf("IsError = false, want true")
	}
	if text := resultText(t, res); !strings.Contains(text, "webmemo login") {
		t.Errorf("text = %q, want it to mention `webmemo login`", text)
	}
}
