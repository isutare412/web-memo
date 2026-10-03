// Package mcpserver exposes the web-memo API as a Model Context Protocol server.
package mcpserver

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/isutare412/web-memo/cli/internal/session"
)

// Options configures the MCP server.
type Options struct {
	// ReadOnly registers only tools that do not modify data.
	ReadOnly bool
	// Version is reported as the server version.
	Version string
}

// tools holds the dependencies shared by tool handlers.
type tools struct {
	session *session.Session
}

// New builds an MCP server whose tools call the API through s.
func New(s *session.Session, opts Options) *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{Name: "webmemo", Version: opts.Version}, nil)

	t := &tools{session: s}
	t.registerRead(srv)
	// Tools that modify data must only be registered when !opts.ReadOnly.

	return srv
}
