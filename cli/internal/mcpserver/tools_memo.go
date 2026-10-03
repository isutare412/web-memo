package mcpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/isutare412/web-memo/cli/internal/apiclient/gen"
	"github.com/isutare412/web-memo/cli/internal/session"
)

// MCP clients assume a non-read-only tool is destructive unless told
// otherwise, so the hint is set explicitly on every write tool. Tools that
// overwrite or remove existing data are destructive; create_memo only adds.
var (
	hintTrue  = true
	hintFalse = false

	additiveAnnotations    = &mcp.ToolAnnotations{DestructiveHint: &hintFalse}
	destructiveAnnotations = &mcp.ToolAnnotations{DestructiveHint: &hintTrue}
)

func (t *tools) registerWrite(srv *mcp.Server) {
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "create_memo",
		Description: "Create a new memo. Returns the created memo (private until published). Tags are free-form names; use search_tags to reuse existing ones.",
		Annotations: additiveAnnotations,
	}, t.createMemo)

	mcp.AddTool(srv, &mcp.Tool{
		Name: "update_memo",
		Description: "Replace fields of a memo. Omitted fields keep their current value; an empty content string or empty tags array clears that field (the title cannot be blank). " +
			"version defaults to the memo's current version; pass the version you read to fail instead of overwriting a concurrent change. " +
			"For small changes to the content prefer edit_memo. Returns the updated memo.",
		Annotations: destructiveAnnotations,
	}, t.updateMemo)

	mcp.AddTool(srv, &mcp.Tool{
		Name: "edit_memo",
		Description: "Edit a memo's content by replacing old_string with new_string (exact match, whitespace included). " +
			"old_string must match exactly once unless replace_all is true. Returns the updated memo.",
		Annotations: destructiveAnnotations,
	}, t.editMemo)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "delete_memo",
		Description: "Permanently delete a memo by id. This cannot be undone.",
		Annotations: destructiveAnnotations,
	}, t.deleteMemo)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "set_memo_tags",
		Description: "Replace all tags of a memo with the given list (an empty list removes every tag). Returns the updated memo.",
		Annotations: destructiveAnnotations,
	}, t.setMemoTags)

	mcp.AddTool(srv, &mcp.Tool{
		Name: "publish_memo",
		Description: "Set a memo's publish state. private: only the owner and collaborators can see it. " +
			"shared: anyone can open a landing page, but the content is visible only to the owner, collaborators, and approved subscribers " +
			"(subscription requests wait for approval). published: anyone can read the content; going from shared to published approves all pending subscriptions. " +
			"WARNING: changing shared or published to private permanently removes all subscribers and collaborators; this cannot be undone. " +
			"Returns the updated memo.",
		Annotations: destructiveAnnotations,
	}, t.publishMemo)
}

type createMemoInput struct {
	Title   string   `json:"title" jsonschema:"memo title (required, non-empty)"`
	Content string   `json:"content,omitempty" jsonschema:"memo body in markdown"`
	Tags    []string `json:"tags,omitempty" jsonschema:"tag names"`
}

func (t *tools) createMemo(ctx context.Context, _ *mcp.CallToolRequest, in createMemoInput) (*mcp.CallToolResult, any, error) {
	req := gen.CreateMemoRequest{Title: in.Title, Tags: in.Tags}
	if in.Content != "" {
		req.Content = &in.Content
	}
	res, err := t.session.API().CreateMemoWithResponse(ctx, req)
	if err != nil {
		return errorResult(err), nil, nil
	}
	if err := session.CheckStatus(res.StatusCode(), res.Body); err != nil {
		return errorResult(err), nil, nil
	}
	if res.JSON200 == nil {
		return errorResult(errUnexpectedBody), nil, nil
	}
	return respond(res.JSON200)
}

type updateMemoInput struct {
	ID            string    `json:"id" jsonschema:"memo id (UUID)"`
	Title         *string   `json:"title,omitempty" jsonschema:"new title; omit to keep"`
	Content       *string   `json:"content,omitempty" jsonschema:"new full content; omit to keep, empty string clears"`
	Tags          *[]string `json:"tags,omitempty" jsonschema:"new full tag list; omit to keep, empty array clears"`
	Version       *int      `json:"version,omitempty" jsonschema:"version the change is based on; omit to use the current version"`
	PinUpdateTime *bool     `json:"pin_update_time,omitempty" jsonschema:"keep the memo's update time unchanged"`
}

func (t *tools) updateMemo(ctx context.Context, _ *mcp.CallToolRequest, in updateMemoInput) (*mcp.CallToolResult, any, error) {
	id, err := parseMemoID(in.ID)
	if err != nil {
		return errorResult(err), nil, nil
	}
	memo, err := t.fetchMemo(ctx, id)
	if err != nil {
		return errorResult(err), nil, nil
	}

	req := replaceRequest(memo)
	if in.Title != nil {
		req.Title = *in.Title
	}
	if in.Content != nil {
		req.Content = in.Content
	}
	if in.Tags != nil {
		req.Tags = *in.Tags
	}
	if in.Version != nil {
		req.Version = *in.Version
	}

	updated, _, err := t.replaceMemo(ctx, id, in.PinUpdateTime, req)
	if err != nil {
		return errorResult(err), nil, nil
	}
	return respond(updated)
}

type editMemoInput struct {
	ID         string `json:"id" jsonschema:"memo id (UUID)"`
	OldString  string `json:"old_string" jsonschema:"exact text to replace; must be unique in the content unless replace_all is set"`
	NewString  string `json:"new_string" jsonschema:"replacement text (may be empty to delete old_string)"`
	ReplaceAll bool   `json:"replace_all,omitempty" jsonschema:"replace every occurrence of old_string"`
}

func (t *tools) editMemo(ctx context.Context, _ *mcp.CallToolRequest, in editMemoInput) (*mcp.CallToolResult, any, error) {
	id, err := parseMemoID(in.ID)
	if err != nil {
		return errorResult(err), nil, nil
	}

	// The memo may change between the read and the write. On a version
	// conflict, re-read and re-apply the edit once, then give up.
	const attempts = 2
	var lastErr error
	for range attempts {
		memo, err := t.fetchMemo(ctx, id)
		if err != nil {
			return errorResult(err), nil, nil
		}
		content, err := applyEdit(memo.Content, in.OldString, in.NewString, in.ReplaceAll)
		if err != nil {
			return errorResult(err), nil, nil
		}

		req := replaceRequest(memo)
		req.Content = &content
		updated, conflict, err := t.replaceMemo(ctx, id, nil, req)
		if err == nil {
			return respond(updated)
		}
		if !conflict {
			return errorResult(err), nil, nil
		}
		lastErr = err
	}
	return errorResult(fmt.Errorf("memo kept changing while editing: %w", lastErr)), nil, nil
}

func (t *tools) deleteMemo(ctx context.Context, _ *mcp.CallToolRequest, in memoIDInput) (*mcp.CallToolResult, any, error) {
	id, err := parseMemoID(in.ID)
	if err != nil {
		return errorResult(err), nil, nil
	}
	res, err := t.session.API().DeleteMemoWithResponse(ctx, id)
	if err != nil {
		return errorResult(err), nil, nil
	}
	if err := session.CheckStatus(res.StatusCode(), res.Body); err != nil {
		return errorResult(err), nil, nil
	}
	return respond(map[string]string{"deleted": id.String()})
}

type setMemoTagsInput struct {
	ID   string   `json:"id" jsonschema:"memo id (UUID)"`
	Tags []string `json:"tags" jsonschema:"the complete new tag list; empty array removes all tags"`
}

func (t *tools) setMemoTags(ctx context.Context, _ *mcp.CallToolRequest, in setMemoTagsInput) (*mcp.CallToolResult, any, error) {
	id, err := parseMemoID(in.ID)
	if err != nil {
		return errorResult(err), nil, nil
	}
	tags := in.Tags
	if tags == nil {
		tags = []string{} // send [] rather than null
	}
	res, err := t.session.API().ReplaceMemoTagsWithResponse(ctx, id, tags)
	if err != nil {
		return errorResult(err), nil, nil
	}
	if err := session.CheckStatus(res.StatusCode(), res.Body); err != nil {
		return errorResult(err), nil, nil
	}

	// The endpoint answers with the tag list only; read the memo back so every
	// write tool returns the full updated memo.
	memo, err := t.fetchMemo(ctx, id)
	if err != nil {
		return errorResult(err), nil, nil
	}
	return respond(memo)
}

type publishMemoInput struct {
	ID    string `json:"id" jsonschema:"memo id (UUID)"`
	State string `json:"state" jsonschema:"private, shared, or published"`
}

func (t *tools) publishMemo(ctx context.Context, _ *mcp.CallToolRequest, in publishMemoInput) (*mcp.CallToolResult, any, error) {
	id, err := parseMemoID(in.ID)
	if err != nil {
		return errorResult(err), nil, nil
	}
	state := gen.PublishState(in.State)
	switch state {
	case gen.PublishStatePrivate, gen.PublishStateShared, gen.PublishStatePublished:
	default:
		return errorResult(fmt.Errorf("invalid state %q: must be private, shared, or published", in.State)), nil, nil
	}

	res, err := t.session.API().PublishMemoWithResponse(ctx, id, gen.PublishMemoRequest{PublishState: state})
	if err != nil {
		return errorResult(err), nil, nil
	}
	if err := session.CheckStatus(res.StatusCode(), res.Body); err != nil {
		return errorResult(err), nil, nil
	}
	if res.JSON200 == nil {
		return errorResult(errUnexpectedBody), nil, nil
	}
	return respond(res.JSON200)
}

// fetchMemo reads a memo, returning API and transport failures as errors.
func (t *tools) fetchMemo(ctx context.Context, id gen.MemoIDPath) (*gen.Memo, error) {
	res, err := t.session.API().GetMemoWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := session.CheckStatus(res.StatusCode(), res.Body); err != nil {
		return nil, err
	}
	if res.JSON200 == nil {
		return nil, errUnexpectedBody
	}
	return res.JSON200, nil
}

// replaceRequest builds a PUT body that leaves memo unchanged.
func replaceRequest(memo *gen.Memo) gen.ReplaceMemoRequest {
	content := memo.Content
	return gen.ReplaceMemoRequest{
		Title:   memo.Title,
		Content: &content,
		Tags:    memo.Tags,
		Version: memo.Version,
	}
}

// replaceMemo PUTs req. conflict reports a version conflict (HTTP 409), which
// the caller may retry after re-reading the memo.
func (t *tools) replaceMemo(
	ctx context.Context, id gen.MemoIDPath, pinUpdateTime *bool, req gen.ReplaceMemoRequest,
) (memo *gen.Memo, conflict bool, err error) {
	res, err := t.session.API().ReplaceMemoWithResponse(ctx, id, &gen.ReplaceMemoParams{PinUpdateTime: pinUpdateTime}, req)
	if err != nil {
		return nil, false, err
	}
	if err := session.CheckStatus(res.StatusCode(), res.Body); err != nil {
		return nil, res.StatusCode() == http.StatusConflict, err
	}
	if res.JSON200 == nil {
		return nil, false, errUnexpectedBody
	}
	return res.JSON200, false, nil
}
