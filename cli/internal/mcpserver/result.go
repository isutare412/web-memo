package mcpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/isutare412/web-memo/cli/internal/apiclient/gen"
)

// previewRunes is the maximum number of runes of memo content kept in a
// search result preview.
const previewRunes = 200

// errUnexpectedBody is reported when the API answers 2xx with a body that is
// not the expected JSON.
var errUnexpectedBody = errors.New("api returned an unexpected response body")

// jsonResult wraps v as indented JSON text content.
func jsonResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode result: %w", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, nil
}

// errorResult reports err to the agent as a tool error (IsError) instead of a
// protocol error, so the model can read the message and react.
func errorResult(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
	}
}

// preview truncates s to previewRunes runes, appending an ellipsis when cut.
// It never splits a multi-byte character.
func preview(s string) string {
	if utf8.RuneCountInString(s) <= previewRunes {
		return s
	}
	return string([]rune(s)[:previewRunes]) + "…"
}

// parseMemoID parses a tool-supplied memo id.
func parseMemoID(s string) (gen.MemoIDPath, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return gen.MemoIDPath{}, fmt.Errorf("invalid memo id %q: must be a UUID", s)
	}
	return id, nil
}

// respond adapts jsonResult to the typed tool handler return shape.
func respond(v any) (*mcp.CallToolResult, any, error) {
	res, err := jsonResult(v)
	return res, nil, err
}
