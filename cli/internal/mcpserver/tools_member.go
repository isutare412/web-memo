package mcpserver

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/isutare412/web-memo/cli/internal/apiclient/gen"
	"github.com/isutare412/web-memo/cli/internal/session"
)

// registerMember registers the subscription and collaboration tools. All of
// them are destructive: besides adding state, each can also remove it
// (unsubscribe, cancel) or revoke an approval (authorize_member).
func (t *tools) registerMember(srv *mcp.Server) {
	mcp.AddTool(srv, &mcp.Tool{
		Name: "subscribe_memo",
		Description: "Subscribe to, or unsubscribe from, another user's memo as the signed-in user. " +
			"action=subscribe: the memo must be shared or published and not your own. Subscribing to a published memo is approved immediately; " +
			"for a shared memo the subscription stays pending until the owner approves it (authorize_member), and the memo content stays hidden until then. " +
			"Returns the subscription including its approved flag. " +
			"action=unsubscribe: drops your subscription (pending or approved) to a memo you do not own.",
		Annotations: destructiveAnnotations,
	}, t.subscribeMemo)

	mcp.AddTool(srv, &mcp.Tool{
		Name: "request_collaboration",
		Description: "Request to collaborate on, or cancel your collaboration on, another user's memo as the signed-in user. " +
			"action=request: you must be able to see the memo (for a shared memo your subscription must already be approved). " +
			"The request stays pending until the owner approves it (authorize_member). " +
			"action=cancel: withdraws your pending request or ends your collaboration.",
		Annotations: destructiveAnnotations,
	}, t.requestCollaboration)

	mcp.AddTool(srv, &mcp.Tool{
		Name: "authorize_member",
		Description: "Memo owner only: approve or reject a subscriber or collaborator of your memo; any other caller gets a permission error. " +
			"Find user ids and pending requests with list_memo_members. " +
			"approve=true grants access; approve=false sets the user's approval to false, which also revokes access that was previously approved " +
			"(the entry stays in the member list as unapproved). " +
			"Fails if the user is already in the requested state.",
		Annotations: destructiveAnnotations,
	}, t.authorizeMember)
}

// currentUserID returns the id of the signed-in user. The member endpoints
// require the caller's own id in the path and reject any other value.
func (t *tools) currentUserID(ctx context.Context) (gen.UserIDPath, error) {
	res, err := t.session.API().GetCurrentUserWithResponse(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("look up current user: %w", err)
	}
	if err := session.CheckStatus(res.StatusCode(), res.Body); err != nil {
		return uuid.Nil, err
	}
	if res.JSON200 == nil {
		return uuid.Nil, errUnexpectedBody
	}
	return res.JSON200.ID, nil
}

func parseUserID(s string) (gen.UserIDPath, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id %q: must be a UUID", s)
	}
	return id, nil
}

func invalidChoice(field, got string, allowed ...string) error {
	return fmt.Errorf("invalid %s %q: must be one of %v", field, got, allowed)
}

// memberAction is the summary returned by tools whose API call has no body.
type memberAction struct {
	MemoID string `json:"memoId"`
	Action string `json:"action"`
}

type subscribeMemoInput struct {
	ID     string `json:"id" jsonschema:"memo id (UUID)"`
	Action string `json:"action" jsonschema:"subscribe or unsubscribe"`
}

func (t *tools) subscribeMemo(ctx context.Context, _ *mcp.CallToolRequest, in subscribeMemoInput) (*mcp.CallToolResult, any, error) {
	if in.Action != "subscribe" && in.Action != "unsubscribe" {
		return errorResult(invalidChoice("action", in.Action, "subscribe", "unsubscribe")), nil, nil
	}
	memoID, err := parseMemoID(in.ID)
	if err != nil {
		return errorResult(err), nil, nil
	}
	userID, err := t.currentUserID(ctx)
	if err != nil {
		return errorResult(err), nil, nil
	}

	if in.Action == "unsubscribe" {
		res, err := t.session.API().UnsubscribeWithResponse(ctx, memoID, userID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		if err := session.CheckStatus(res.StatusCode(), res.Body); err != nil {
			return errorResult(err), nil, nil
		}
		return respond(memberAction{MemoID: in.ID, Action: in.Action})
	}

	res, err := t.session.API().SubscribeWithResponse(ctx, memoID, userID)
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

type requestCollaborationInput struct {
	ID     string `json:"id" jsonschema:"memo id (UUID)"`
	Action string `json:"action" jsonschema:"request or cancel"`
}

func (t *tools) requestCollaboration(ctx context.Context, _ *mcp.CallToolRequest, in requestCollaborationInput) (*mcp.CallToolResult, any, error) {
	if in.Action != "request" && in.Action != "cancel" {
		return errorResult(invalidChoice("action", in.Action, "request", "cancel")), nil, nil
	}
	memoID, err := parseMemoID(in.ID)
	if err != nil {
		return errorResult(err), nil, nil
	}
	userID, err := t.currentUserID(ctx)
	if err != nil {
		return errorResult(err), nil, nil
	}

	var (
		status int
		body   []byte
	)
	if in.Action == "request" {
		res, err := t.session.API().RequestCollaborationWithResponse(ctx, memoID, userID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		status, body = res.StatusCode(), res.Body
	} else {
		res, err := t.session.API().CancelCollaborationWithResponse(ctx, memoID, userID)
		if err != nil {
			return errorResult(err), nil, nil
		}
		status, body = res.StatusCode(), res.Body
	}
	if err := session.CheckStatus(status, body); err != nil {
		return errorResult(err), nil, nil
	}
	return respond(memberAction{MemoID: in.ID, Action: in.Action})
}

type authorizeMemberInput struct {
	ID      string `json:"id" jsonschema:"memo id (UUID); you must own the memo"`
	UserID  string `json:"user_id" jsonschema:"id (UUID) of the subscriber or collaborator to decide on"`
	Kind    string `json:"kind" jsonschema:"subscriber or collaborator"`
	Approve bool   `json:"approve" jsonschema:"true to approve, false to reject or revoke"`
}

type authorizeMemberResult struct {
	MemoID   string `json:"memoId"`
	UserID   string `json:"userId"`
	Kind     string `json:"kind"`
	Approved bool   `json:"approved"`
}

func (t *tools) authorizeMember(ctx context.Context, _ *mcp.CallToolRequest, in authorizeMemberInput) (*mcp.CallToolResult, any, error) {
	if in.Kind != "subscriber" && in.Kind != "collaborator" {
		return errorResult(invalidChoice("kind", in.Kind, "subscriber", "collaborator")), nil, nil
	}
	memoID, err := parseMemoID(in.ID)
	if err != nil {
		return errorResult(err), nil, nil
	}
	userID, err := parseUserID(in.UserID)
	if err != nil {
		return errorResult(err), nil, nil
	}

	var (
		status int
		body   []byte
	)
	if in.Kind == "subscriber" {
		res, err := t.session.API().AuthorizeSubscriptionWithResponse(ctx, memoID, userID, gen.AuthorizeSubscriptionRequest{Approve: in.Approve})
		if err != nil {
			return errorResult(err), nil, nil
		}
		status, body = res.StatusCode(), res.Body
	} else {
		res, err := t.session.API().AuthorizeCollaborationWithResponse(ctx, memoID, userID, gen.AuthorizeCollaborationRequest{Approve: in.Approve})
		if err != nil {
			return errorResult(err), nil, nil
		}
		status, body = res.StatusCode(), res.Body
	}
	if err := session.CheckStatus(status, body); err != nil {
		return errorResult(err), nil, nil
	}
	return respond(authorizeMemberResult{MemoID: in.ID, UserID: in.UserID, Kind: in.Kind, Approved: in.Approve})
}
