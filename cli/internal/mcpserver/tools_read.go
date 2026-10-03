package mcpserver

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/isutare412/web-memo/cli/internal/apiclient/gen"
	"github.com/isutare412/web-memo/cli/internal/session"
)

var readOnlyAnnotations = &mcp.ToolAnnotations{ReadOnlyHint: true}

func (t *tools) registerRead(srv *mcp.Server) {
	mcp.AddTool(srv, &mcp.Tool{
		Name: "search_memos",
		Description: "List or search the user's memos. With query, runs a hybrid semantic + keyword search " +
			"(results include relevance scores); without it, lists memos by sort order. " +
			"tags keeps only memos that have the given tags. sort is createTime or updateTime (default updateTime). " +
			"Results are paged (page starts at 1) and contain a content preview, not the full body; use get_memo for the full content.",
		Annotations: readOnlyAnnotations,
	}, t.searchMemos)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get_memo",
		Description: "Get one memo by id, including its full content, tags, publish state, and version.",
		Annotations: readOnlyAnnotations,
	}, t.getMemo)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list_memo_members",
		Description: "List the subscribers and collaborators of a memo by memo id, with their approval state.",
		Annotations: readOnlyAnnotations,
	}, t.listMemoMembers)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "search_tags",
		Description: "Search the user's existing tags by keyword. Returns an array of matching tag names.",
		Annotations: readOnlyAnnotations,
	}, t.searchTags)
}

type searchMemosInput struct {
	Query    string   `json:"query,omitempty" jsonschema:"search text; omit to list memos without searching"`
	Tags     []string `json:"tags,omitempty" jsonschema:"only memos having these tags"`
	Sort     string   `json:"sort,omitempty" jsonschema:"createTime or updateTime (default updateTime)"`
	Page     int      `json:"page,omitempty" jsonschema:"page number starting at 1"`
	PageSize int      `json:"page_size,omitempty" jsonschema:"number of memos per page"`
}

type memoSummary struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	Tags         []string         `json:"tags"`
	PublishState gen.PublishState `json:"publishState"`
	CreateTime   time.Time        `json:"createTime"`
	UpdateTime   time.Time        `json:"updateTime"`
	Preview      string           `json:"preview"`
	Scores       *gen.MemoScores  `json:"scores,omitempty"`
}

type searchMemosOutput struct {
	Page           *int          `json:"page"`
	PageSize       *int          `json:"pageSize"`
	LastPage       *int          `json:"lastPage"`
	TotalMemoCount *int          `json:"totalMemoCount"`
	Memos          []memoSummary `json:"memos"`
}

func (t *tools) searchMemos(ctx context.Context, _ *mcp.CallToolRequest, in searchMemosInput) (*mcp.CallToolResult, any, error) {
	params := &gen.ListMemosParams{Tag: in.Tags}
	if in.Query != "" {
		params.SearchQuery = &in.Query
	}
	if in.Sort != "" {
		sort := gen.MemoSortKey(in.Sort)
		if sort != gen.MemoSortKeyCreateTime && sort != gen.MemoSortKeyUpdateTime {
			return errorResult(fmt.Errorf("invalid sort %q: must be createTime or updateTime", in.Sort)), nil, nil
		}
		params.Sort = &sort
	}
	if in.Page > 0 {
		params.Page = &in.Page
	}
	if in.PageSize > 0 {
		params.PageSize = &in.PageSize
	}

	res, err := t.session.API().ListMemosWithResponse(ctx, params)
	if err != nil {
		return errorResult(err), nil, nil
	}
	if err := session.CheckStatus(res.StatusCode(), res.Body); err != nil {
		return errorResult(err), nil, nil
	}

	if res.JSON200 == nil {
		return errorResult(errUnexpectedBody), nil, nil
	}
	out := searchMemosOutput{
		Page:           res.JSON200.Page,
		PageSize:       res.JSON200.PageSize,
		LastPage:       res.JSON200.LastPage,
		TotalMemoCount: res.JSON200.TotalMemoCount,
		Memos:          make([]memoSummary, 0, len(res.JSON200.Memos)),
	}
	for _, m := range res.JSON200.Memos {
		tags := m.Tags
		if tags == nil {
			tags = []string{}
		}
		out.Memos = append(out.Memos, memoSummary{
			ID:           m.ID.String(),
			Title:        m.Title,
			Tags:         tags,
			PublishState: m.PublishState,
			CreateTime:   m.CreateTime,
			UpdateTime:   m.UpdateTime,
			Preview:      preview(m.Content),
			Scores:       m.Scores,
		})
	}
	return respond(out)
}

type memoIDInput struct {
	ID string `json:"id" jsonschema:"memo id (UUID)"`
}

func (t *tools) getMemo(ctx context.Context, _ *mcp.CallToolRequest, in memoIDInput) (*mcp.CallToolResult, any, error) {
	id, err := parseMemoID(in.ID)
	if err != nil {
		return errorResult(err), nil, nil
	}

	res, err := t.session.API().GetMemoWithResponse(ctx, id)
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

type memoMembersOutput struct {
	MemoOwnerID   uuid.UUID          `json:"memoOwnerId"`
	Subscribers   []gen.Subscriber   `json:"subscribers"`
	Collaborators []gen.Collaborator `json:"collaborators"`
}

func (t *tools) listMemoMembers(ctx context.Context, _ *mcp.CallToolRequest, in memoIDInput) (*mcp.CallToolResult, any, error) {
	id, err := parseMemoID(in.ID)
	if err != nil {
		return errorResult(err), nil, nil
	}

	subs, err := t.session.API().ListSubscribersWithResponse(ctx, id)
	if err != nil {
		return errorResult(err), nil, nil
	}
	if err := session.CheckStatus(subs.StatusCode(), subs.Body); err != nil {
		return errorResult(err), nil, nil
	}

	collabs, err := t.session.API().ListCollaboratorsWithResponse(ctx, id)
	if err != nil {
		return errorResult(err), nil, nil
	}
	if err := session.CheckStatus(collabs.StatusCode(), collabs.Body); err != nil {
		return errorResult(err), nil, nil
	}

	if subs.JSON200 == nil || collabs.JSON200 == nil {
		return errorResult(errUnexpectedBody), nil, nil
	}
	out := memoMembersOutput{
		MemoOwnerID:   subs.JSON200.MemoOwnerID,
		Subscribers:   subs.JSON200.Subscribers,
		Collaborators: collabs.JSON200.Collaborators,
	}
	if out.Subscribers == nil {
		out.Subscribers = []gen.Subscriber{}
	}
	if out.Collaborators == nil {
		out.Collaborators = []gen.Collaborator{}
	}
	return respond(out)
}

type searchTagsInput struct {
	Keyword string `json:"keyword" jsonschema:"keyword to match tag names against"`
}

func (t *tools) searchTags(ctx context.Context, _ *mcp.CallToolRequest, in searchTagsInput) (*mcp.CallToolResult, any, error) {
	res, err := t.session.API().SearchTagsWithResponse(ctx, &gen.SearchTagsParams{Kw: &in.Keyword})
	if err != nil {
		return errorResult(err), nil, nil
	}
	if err := session.CheckStatus(res.StatusCode(), res.Body); err != nil {
		return errorResult(err), nil, nil
	}
	if res.JSON200 == nil {
		return errorResult(errUnexpectedBody), nil, nil
	}
	return respond(*res.JSON200)
}
