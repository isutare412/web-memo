package mcpserver

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

const (
	usersMePath = "/api/v1/users/me"
	meJSON      = `{"id":"` + testUserID + `","email":"me@example.com","userName":"me","userType":"user",` +
		`"issuedAt":"2026-01-01T00:00:00Z","expireAt":"2026-01-02T00:00:00Z"}`
)

func TestSubscribeMemo(t *testing.T) {
	subscribers := memoPath + "/subscribers/" + testUserID
	subscribedBody := `{"subscription":{"approved":false,"memoId":"` + testMemoID + `","userId":"` + testUserID + `"}}`

	for _, tc := range []struct {
		name, action, method, path string
		status                     int
		body                       string
	}{
		{"subscribe", "subscribe", "PUT", subscribers, 200, subscribedBody},
		{"unsubscribe", "unsubscribe", "DELETE", subscribers, 200, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api := newFakeMemoAPI(t).
				on("GET", usersMePath, 200, meJSON).
				on(tc.method, tc.path, tc.status, tc.body)
			cs := newTestClient(t, api, Options{})

			res := callTool(t, cs, "subscribe_memo", map[string]any{"id": testMemoID, "action": tc.action})
			if res.IsError {
				t.Fatalf("unexpected error: %s", resultText(t, res))
			}
			if len(api.requests) != 2 || len(api.calls(tc.method, tc.path)) != 1 {
				t.Fatalf("requests = %+v, want GET users/me then %s %s", api.requests, tc.method, tc.path)
			}
			var out map[string]any
			if err := json.Unmarshal([]byte(resultText(t, res)), &out); err != nil {
				t.Fatalf("result is not JSON: %v", err)
			}
			if tc.action == "subscribe" {
				sub, _ := out["subscription"].(map[string]any)
				if sub["approved"] != false || sub["memoId"] != testMemoID {
					t.Errorf("subscribe result = %v, want the subscription with approved=false", out)
				}
			} else if out["action"] != "unsubscribe" || out["memoId"] != testMemoID {
				t.Errorf("unsubscribe result = %v", out)
			}
		})
	}
}

func TestRequestCollaboration(t *testing.T) {
	collaborator := memoPath + "/collaborators/" + testUserID
	for _, tc := range []struct{ action, method string }{
		{"request", "PUT"},
		{"cancel", "DELETE"},
	} {
		t.Run(tc.action, func(t *testing.T) {
			api := newFakeMemoAPI(t).
				on("GET", usersMePath, 200, meJSON).
				on(tc.method, collaborator, 200, "")
			cs := newTestClient(t, api, Options{})

			res := callTool(t, cs, "request_collaboration", map[string]any{"id": testMemoID, "action": tc.action})
			if res.IsError {
				t.Fatalf("unexpected error: %s", resultText(t, res))
			}
			if len(api.requests) != 2 || len(api.calls(tc.method, collaborator)) != 1 {
				t.Fatalf("requests = %+v, want GET users/me then %s %s", api.requests, tc.method, collaborator)
			}
			var out map[string]any
			if err := json.Unmarshal([]byte(resultText(t, res)), &out); err != nil {
				t.Fatalf("result is not JSON: %v", err)
			}
			if out["action"] != tc.action || out["memoId"] != testMemoID {
				t.Errorf("result = %v", out)
			}
		})
	}
}

func TestAuthorizeMember(t *testing.T) {
	for _, tc := range []struct {
		kind, segment string
		approve       bool
	}{
		{"subscriber", "subscribers", true},
		{"subscriber", "subscribers", false},
		{"collaborator", "collaborators", true},
		{"collaborator", "collaborators", false},
	} {
		name := tc.kind + "/approve=" + map[bool]string{true: "true", false: "false"}[tc.approve]
		t.Run(name, func(t *testing.T) {
			path := memoPath + "/" + tc.segment + "/" + testUserID + "/authorize"
			api := newFakeMemoAPI(t).on("POST", path, 200, "")
			cs := newTestClient(t, api, Options{})

			res := callTool(t, cs, "authorize_member", map[string]any{
				"id": testMemoID, "user_id": testUserID, "kind": tc.kind, "approve": tc.approve,
			})
			if res.IsError {
				t.Fatalf("unexpected error: %s", resultText(t, res))
			}
			calls := api.calls("POST", path)
			if len(api.requests) != 1 || len(calls) != 1 {
				t.Fatalf("requests = %+v, want exactly POST %s", api.requests, path)
			}
			if len(calls[0].Body) != 1 || calls[0].Body["approve"] != tc.approve {
				t.Errorf("body = %v, want {approve: %v}", calls[0].Body, tc.approve)
			}
			var out map[string]any
			if err := json.Unmarshal([]byte(resultText(t, res)), &out); err != nil {
				t.Fatalf("result is not JSON: %v", err)
			}
			if out["kind"] != tc.kind || out["userId"] != testUserID || out["approved"] != tc.approve {
				t.Errorf("result = %v", out)
			}
		})
	}
}

func TestAuthorizeMemberRequiresApprove(t *testing.T) {
	api := newFakeMemoAPI(t)
	cs := newTestClient(t, api, Options{})
	res := callTool(t, cs, "authorize_member", map[string]any{"id": testMemoID, "user_id": testUserID, "kind": "subscriber"})
	if !res.IsError {
		t.Errorf("want error when approve is omitted, got %s", resultText(t, res))
	}
	if len(api.requests) != 0 {
		t.Errorf("api called: %+v", api.requests)
	}
}

func TestMemberToolsRejectInvalidInput(t *testing.T) {
	for _, tc := range []struct {
		name, tool string
		args       map[string]any
		wantMsg    string
	}{
		{"subscribe bad action", "subscribe_memo", map[string]any{"id": testMemoID, "action": "follow"}, "action"},
		{"subscribe empty action", "subscribe_memo", map[string]any{"id": testMemoID, "action": ""}, "action"},
		{"subscribe bad id", "subscribe_memo", map[string]any{"id": "nope", "action": "subscribe"}, "invalid memo id"},
		{"collab bad action", "request_collaboration", map[string]any{"id": testMemoID, "action": "subscribe"}, "action"},
		{"collab bad id", "request_collaboration", map[string]any{"id": "nope", "action": "request"}, "invalid memo id"},
		{"authorize bad kind", "authorize_member", map[string]any{"id": testMemoID, "user_id": testUserID, "kind": "owner", "approve": true}, "kind"},
		{"authorize bad id", "authorize_member", map[string]any{"id": "nope", "user_id": testUserID, "kind": "subscriber", "approve": true}, "invalid memo id"},
		{"authorize bad user id", "authorize_member", map[string]any{"id": testMemoID, "user_id": "nope", "kind": "subscriber", "approve": true}, "invalid user id"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api := newFakeMemoAPI(t)
			cs := newTestClient(t, api, Options{})
			res := callTool(t, cs, tc.tool, tc.args)
			if !res.IsError {
				t.Fatalf("want tool error, got %s", resultText(t, res))
			}
			if msg := resultText(t, res); !strings.Contains(msg, tc.wantMsg) {
				t.Errorf("error %q does not mention %q", msg, tc.wantMsg)
			}
			if len(api.requests) != 0 {
				t.Errorf("api called: %+v", api.requests)
			}
		})
	}
}

func TestMemberToolsReportApiErrors(t *testing.T) {
	t.Run("subscribe rejected", func(t *testing.T) {
		api := newFakeMemoAPI(t).
			on("GET", usersMePath, 200, meJSON).
			on("PUT", memoPath+"/subscribers/"+testUserID, 400, `{"msg":"cannot subscribe memo of your own"}`)
		cs := newTestClient(t, api, Options{})
		res := callTool(t, cs, "subscribe_memo", map[string]any{"id": testMemoID, "action": "subscribe"})
		if !res.IsError || !strings.Contains(resultText(t, res), "cannot subscribe memo of your own") {
			t.Errorf("want API message as tool error, got %+v", res)
		}
	})
	t.Run("current user lookup fails", func(t *testing.T) {
		api := newFakeMemoAPI(t).on("GET", usersMePath, 500, `{"msg":"boom"}`)
		cs := newTestClient(t, api, Options{})
		res := callTool(t, cs, "request_collaboration", map[string]any{"id": testMemoID, "action": "request"})
		if !res.IsError {
			t.Errorf("want tool error, got %s", resultText(t, res))
		}
		if len(api.requests) != 1 {
			t.Errorf("requests = %+v, want only the users/me lookup", api.requests)
		}
	})
	t.Run("authorize denied", func(t *testing.T) {
		api := newFakeMemoAPI(t).
			on("POST", memoPath+"/collaborators/"+testUserID+"/authorize", 403, `{"msg":"not allowed to access memo"}`)
		cs := newTestClient(t, api, Options{})
		res := callTool(t, cs, "authorize_member", map[string]any{"id": testMemoID, "user_id": testUserID, "kind": "collaborator", "approve": true})
		if !res.IsError || !strings.Contains(resultText(t, res), "not allowed") {
			t.Errorf("want API message as tool error, got %+v", res)
		}
	})
}

func TestMemberToolsHiddenWhenReadOnly(t *testing.T) {
	cs := newTestClient(t, http.NotFoundHandler(), Options{ReadOnly: true})
	for _, name := range toolNames(t, cs) {
		switch name {
		case "subscribe_memo", "request_collaboration", "authorize_member":
			t.Errorf("%s registered in read-only mode", name)
		}
	}
}
