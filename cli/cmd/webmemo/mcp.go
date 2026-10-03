package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/isutare412/web-memo/cli/internal/credential"
	"github.com/isutare412/web-memo/cli/internal/mcpserver"
	"github.com/isutare412/web-memo/cli/internal/session"
)

// mcpHTTPTimeout is longer than httpTimeout because image uploads and
// waiting for image processing can take a while.
const mcpHTTPTimeout = 2 * time.Minute

// runMCP serves the MCP protocol on stdin/stdout until the client disconnects
// or ctx is canceled. stdout carries only JSON-RPC; the only other output is
// the -h usage text written to out.
func runMCP(ctx context.Context, args []string, env func(string) string, store *credential.Store, out io.Writer) error {
	srv, err := newMCPServer(args, env, store)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printUsage(mcpFlagSet(new(bool)), out)
			return nil
		}
		return err
	}
	return srv.Run(ctx, &mcp.StdioTransport{})
}

func mcpFlagSet(readOnly *bool) *flag.FlagSet {
	fs := newFlagSet("mcp")
	fs.BoolVar(readOnly, "read-only", false, "expose only the read tools (search_memos, get_memo, list_memo_members, search_tags)")
	return fs
}

// newMCPServer parses args and builds the server. A missing token is not an
// error: the server starts and tool calls return login guidance.
func newMCPServer(args []string, env func(string) string, store *credential.Store) (*mcp.Server, error) {
	var readOnly bool
	if err := parseFlags(mcpFlagSet(&readOnly), args); err != nil {
		return nil, err
	}

	cred, src, err := credential.Resolve(store, env)
	if err != nil {
		return nil, err
	}
	sess := session.New(cred, src, store, &http.Client{Timeout: mcpHTTPTimeout})
	return mcpserver.New(sess, mcpserver.Options{ReadOnly: readOnly, Version: buildVersion()}), nil
}

// buildVersion returns the main module version, or "dev" for local builds.
func buildVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}
