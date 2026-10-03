package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
)

var writeTools = []string{"create_memo", "delete_memo", "edit_memo", "publish_memo", "set_memo_tags", "update_memo"}

func memoJSON(title, content string, tags []string, version int, state string) string {
	b, _ := json.Marshal(map[string]any{
		"id": testMemoID, "ownerId": testOwnerID, "title": title, "content": content,
		"tags": tags, "publishState": state, "version": version,
		"createTime": "2026-01-01T00:00:00Z", "updateTime": "2026-01-02T00:00:00Z",
	})
	return string(b)
}

// recorded is a request seen by fakeMemoAPI.
type recorded struct {
	Method, Path, Query string
	Body                map[string]any
}

// fakeMemoAPI answers each request with the next response registered for its
// "METHOD path" key and records all requests.
type fakeMemoAPI struct {
	t         *testing.T
	mu        sync.Mutex
	responses map[string][]fakeResponse
	requests  []recorded
}

type fakeResponse struct {
	status int
	body   string
}

func newFakeMemoAPI(t *testing.T) *fakeMemoAPI {
	return &fakeMemoAPI{t: t, responses: make(map[string][]fakeResponse)}
}

func (f *fakeMemoAPI) on(method, path string, status int, body string) *fakeMemoAPI {
	key := method + " " + path
	f.responses[key] = append(f.responses[key], fakeResponse{status, body})
	return f
}

func (f *fakeMemoAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	rec := recorded{Method: r.Method, Path: r.URL.Path, Query: r.URL.RawQuery}
	if raw, _ := io.ReadAll(r.Body); len(raw) > 0 {
		var decoded any
		if err := json.Unmarshal(raw, &decoded); err != nil {
			f.t.Errorf("request body is not JSON: %q", raw)
		}
		if m, ok := decoded.(map[string]any); ok {
			rec.Body = m
		} else {
			rec.Body = map[string]any{"_": decoded}
		}
	}
	f.requests = append(f.requests, rec)

	key := r.Method + " " + r.URL.Path
	queue := f.responses[key]
	if len(queue) == 0 {
		f.t.Errorf("unexpected request %s", key)
		http.Error(w, "unexpected", http.StatusInternalServerError)
		return
	}
	resp := queue[0]
	if len(queue) > 1 {
		f.responses[key] = queue[1:]
	}
	if resp.body == "" {
		w.WriteHeader(resp.status)
		return
	}
	writeJSON(w, resp.status, resp.body)
}

// calls returns the recorded requests with method to path.
func (f *fakeMemoAPI) calls(method, path string) []recorded {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []recorded
	for _, r := range f.requests {
		if r.Method == method && r.Path == path {
			out = append(out, r)
		}
	}
	return out
}

const memoPath = "/api/v1/memos/" + testMemoID

func TestToolListWriteTools(t *testing.T) {
	cs := newTestClient(t, http.NotFoundHandler(), Options{})
	got := make(map[string]bool)
	for _, name := range toolNames(t, cs) {
		got[name] = true
	}
	for _, name := range append(append([]string{}, readTools...), writeTools...) {
		if !got[name] {
			t.Errorf("tool %q missing from %v", name, got)
		}
	}
}

func TestCreateMemo(t *testing.T) {
	api := newFakeMemoAPI(t).on("POST", "/api/v1/memos", 200, memoJSON("T", "body", []string{"a"}, 1, "private"))
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "create_memo", map[string]any{"title": "T", "content": "body", "tags": []string{"a"}})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	var memo struct {
		ID      string
		Title   string
		Version int
	}
	if err := json.Unmarshal([]byte(resultText(t, res)), &memo); err != nil {
		t.Fatal(err)
	}
	if memo.ID != testMemoID || memo.Title != "T" || memo.Version != 1 {
		t.Errorf("memo = %+v", memo)
	}
	body := api.calls("POST", "/api/v1/memos")[0].Body
	if body["title"] != "T" || body["content"] != "body" {
		t.Errorf("body = %v", body)
	}
	if tags, _ := body["tags"].([]any); len(tags) != 1 || tags[0] != "a" {
		t.Errorf("tags = %v", body["tags"])
	}
}

func TestCreateMemoTitleOnly(t *testing.T) {
	api := newFakeMemoAPI(t).on("POST", "/api/v1/memos", 200, memoJSON("T", "", nil, 1, "private"))
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "create_memo", map[string]any{"title": "T"})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	body := api.calls("POST", "/api/v1/memos")[0].Body
	if _, ok := body["content"]; ok {
		t.Errorf("content should be omitted, body = %v", body)
	}
}

func TestUpdateMemoMergesFields(t *testing.T) {
	existing := memoJSON("old title", "old content", []string{"x", "y"}, 5, "private")
	updated := memoJSON("old title", "new content", []string{"x", "y"}, 6, "private")

	tests := []struct {
		name        string
		args        map[string]any
		wantTitle   string
		wantContent string
		wantTags    []any
		wantVersion float64
		wantQuery   string
	}{
		{
			name: "content only", args: map[string]any{"id": testMemoID, "content": "new content"},
			wantTitle: "old title", wantContent: "new content", wantTags: []any{"x", "y"}, wantVersion: 5,
		},
		{
			name: "explicit version", args: map[string]any{"id": testMemoID, "content": "new content", "version": 3},
			wantTitle: "old title", wantContent: "new content", wantTags: []any{"x", "y"}, wantVersion: 3,
		},
		{
			name: "empty tags", args: map[string]any{"id": testMemoID, "tags": []string{}},
			wantTitle: "old title", wantContent: "old content", wantTags: nil, wantVersion: 5,
		},
		{
			name: "title and pin", args: map[string]any{"id": testMemoID, "title": "new title", "pin_update_time": true},
			wantTitle: "new title", wantContent: "old content", wantTags: []any{"x", "y"}, wantVersion: 5,
			wantQuery: "pinUpdateTime=true",
		},
		{
			name: "empty content", args: map[string]any{"id": testMemoID, "content": ""},
			wantTitle: "old title", wantContent: "", wantTags: []any{"x", "y"}, wantVersion: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newFakeMemoAPI(t).on("GET", memoPath, 200, existing).on("PUT", memoPath, 200, updated)
			cs := newTestClient(t, api, Options{})

			res := callTool(t, cs, "update_memo", tt.args)
			if res.IsError {
				t.Fatalf("unexpected error: %s", resultText(t, res))
			}
			puts := api.calls("PUT", memoPath)
			if len(puts) != 1 {
				t.Fatalf("PUT count = %d", len(puts))
			}
			body := puts[0].Body
			if body["title"] != tt.wantTitle || body["content"] != tt.wantContent || body["version"] != tt.wantVersion {
				t.Errorf("body = %v", body)
			}
			tags, _ := body["tags"].([]any)
			if fmt.Sprint(tags) != fmt.Sprint(tt.wantTags) {
				t.Errorf("tags = %v, want %v", tags, tt.wantTags)
			}
			if puts[0].Query != tt.wantQuery {
				t.Errorf("query = %q, want %q", puts[0].Query, tt.wantQuery)
			}
		})
	}
}

func TestUpdateMemoReportsConflict(t *testing.T) {
	api := newFakeMemoAPI(t).
		on("GET", memoPath, 200, memoJSON("t", "c", nil, 5, "private")).
		on("PUT", memoPath, 409, `{"msg":"memo version is outdated"}`)
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "update_memo", map[string]any{"id": testMemoID, "content": "x", "version": 1})
	if !res.IsError || !strings.Contains(resultText(t, res), "outdated") {
		t.Errorf("result = %+v", res)
	}
}

func TestEditMemo(t *testing.T) {
	api := newFakeMemoAPI(t).
		on("GET", memoPath, 200, memoJSON("t", "hello world", []string{"x"}, 2, "private")).
		on("PUT", memoPath, 200, memoJSON("t", "hello there", []string{"x"}, 3, "private"))
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "edit_memo", map[string]any{"id": testMemoID, "old_string": "world", "new_string": "there"})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	if !strings.Contains(resultText(t, res), "hello there") {
		t.Errorf("result = %s", resultText(t, res))
	}
	body := api.calls("PUT", memoPath)[0].Body
	if body["content"] != "hello there" || body["version"] != float64(2) || body["title"] != "t" {
		t.Errorf("body = %v", body)
	}
}

func TestEditMemoNoMatchDoesNotWrite(t *testing.T) {
	api := newFakeMemoAPI(t).on("GET", memoPath, 200, memoJSON("t", "hello", nil, 2, "private"))
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "edit_memo", map[string]any{"id": testMemoID, "old_string": "zzz", "new_string": "a"})
	if !res.IsError || !strings.Contains(resultText(t, res), "0 matches") {
		t.Errorf("result = %+v", res)
	}
	if n := len(api.calls("PUT", memoPath)); n != 0 {
		t.Errorf("PUT count = %d, want 0", n)
	}
}

func TestEditMemoRetriesOnConflict(t *testing.T) {
	api := newFakeMemoAPI(t).
		on("GET", memoPath, 200, memoJSON("t", "hello world", nil, 2, "private")).
		on("GET", memoPath, 200, memoJSON("t", "hello world!", nil, 4, "private")).
		on("PUT", memoPath, 409, `{"msg":"memo version is outdated"}`).
		on("PUT", memoPath, 200, memoJSON("t", "hello there!", nil, 5, "private"))
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "edit_memo", map[string]any{"id": testMemoID, "old_string": "world", "new_string": "there"})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	puts := api.calls("PUT", memoPath)
	if len(puts) != 2 {
		t.Fatalf("PUT count = %d, want 2", len(puts))
	}
	if puts[0].Body["version"] != float64(2) {
		t.Errorf("first version = %v", puts[0].Body["version"])
	}
	if puts[1].Body["version"] != float64(4) || puts[1].Body["content"] != "hello there!" {
		t.Errorf("second body = %v", puts[1].Body)
	}
}

func TestEditMemoGivesUpAfterSecondConflict(t *testing.T) {
	api := newFakeMemoAPI(t).
		on("GET", memoPath, 200, memoJSON("t", "hello world", nil, 2, "private")).
		on("GET", memoPath, 200, memoJSON("t", "hello world", nil, 3, "private")).
		on("PUT", memoPath, 409, `{"msg":"memo version is outdated"}`)
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "edit_memo", map[string]any{"id": testMemoID, "old_string": "world", "new_string": "there"})
	if !res.IsError || !strings.Contains(resultText(t, res), "outdated") {
		t.Errorf("result = %+v", res)
	}
	if n := len(api.calls("PUT", memoPath)); n != 2 {
		t.Errorf("PUT count = %d, want 2", n)
	}
}

func TestEditMemoApiRejects(t *testing.T) {
	api := newFakeMemoAPI(t).
		on("GET", memoPath, 200, memoJSON("t", "hello world", nil, 2, "private")).
		on("PUT", memoPath, 400, `{"msg":"content too long"}`)
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "edit_memo", map[string]any{"id": testMemoID, "old_string": "world", "new_string": "there"})
	if !res.IsError || !strings.Contains(resultText(t, res), "content too long") {
		t.Errorf("result = %+v", res)
	}
	if n := len(api.calls("PUT", memoPath)); n != 1 {
		t.Errorf("PUT count = %d, want 1 (no retry on 400)", n)
	}
}

func TestDeleteMemo(t *testing.T) {
	api := newFakeMemoAPI(t).on("DELETE", memoPath, 204, "")
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "delete_memo", map[string]any{"id": testMemoID})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	var out map[string]string
	if err := json.Unmarshal([]byte(resultText(t, res)), &out); err != nil {
		t.Fatal(err)
	}
	if out["deleted"] != testMemoID {
		t.Errorf("out = %v", out)
	}
}

func TestDeleteMemoAnnotation(t *testing.T) {
	cs := newTestClient(t, http.NotFoundHandler(), Options{})
	res, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range res.Tools {
		destructive := tool.Annotations != nil && tool.Annotations.DestructiveHint != nil && *tool.Annotations.DestructiveHint
		if tool.Name == "delete_memo" && !destructive {
			t.Errorf("delete_memo DestructiveHint is not true: %+v", tool.Annotations)
		}
		if tool.Name != "delete_memo" && destructive {
			t.Errorf("%s must not be destructive", tool.Name)
		}
	}
}

func TestSetMemoTags(t *testing.T) {
	api := newFakeMemoAPI(t).
		on("PUT", memoPath+"/tags", 200, `["a","b"]`).
		on("GET", memoPath, 200, memoJSON("t", "c", []string{"a", "b"}, 2, "private"))
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "set_memo_tags", map[string]any{"id": testMemoID, "tags": []string{"a", "b"}})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	var memo struct{ Tags []string }
	if err := json.Unmarshal([]byte(resultText(t, res)), &memo); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(memo.Tags) != "[a b]" {
		t.Errorf("tags = %v", memo.Tags)
	}
	body := api.calls("PUT", memoPath+"/tags")[0].Body
	if fmt.Sprint(body["_"]) != "[a b]" {
		t.Errorf("body = %v", body)
	}
}

func TestSetMemoTagsEmptyClears(t *testing.T) {
	api := newFakeMemoAPI(t).
		on("PUT", memoPath+"/tags", 200, `[]`).
		on("GET", memoPath, 200, memoJSON("t", "c", []string{}, 2, "private"))
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "set_memo_tags", map[string]any{"id": testMemoID, "tags": []string{}})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	body := api.calls("PUT", memoPath+"/tags")[0].Body
	if arr, ok := body["_"].([]any); !ok || len(arr) != 0 {
		t.Errorf("body = %v, want empty array", body)
	}
}

func TestPublishMemo(t *testing.T) {
	api := newFakeMemoAPI(t).on("POST", memoPath+"/publish", 200, memoJSON("t", "c", nil, 2, "published"))
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "publish_memo", map[string]any{"id": testMemoID, "state": "published"})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	if !strings.Contains(resultText(t, res), `"publishState": "published"`) {
		t.Errorf("result = %s", resultText(t, res))
	}
	if body := api.calls("POST", memoPath+"/publish")[0].Body; body["publishState"] != "published" {
		t.Errorf("body = %v", body)
	}
}

func TestPublishMemoRejectsUnknownState(t *testing.T) {
	api := newFakeMemoAPI(t)
	cs := newTestClient(t, api, Options{})

	res := callTool(t, cs, "publish_memo", map[string]any{"id": testMemoID, "state": "public"})
	if !res.IsError || !strings.Contains(resultText(t, res), "private") {
		t.Errorf("result = %+v", res)
	}
	if len(api.requests) != 0 {
		t.Errorf("api was called: %v", api.requests)
	}
}

func TestWriteToolsInvalidID(t *testing.T) {
	cs := newTestClient(t, newFakeMemoAPI(t), Options{})
	for name, args := range map[string]map[string]any{
		"update_memo":   {"id": "nope", "title": "t"},
		"edit_memo":     {"id": "nope", "old_string": "a", "new_string": "b"},
		"delete_memo":   {"id": "nope"},
		"set_memo_tags": {"id": "nope", "tags": []string{}},
		"publish_memo":  {"id": "nope", "state": "private"},
	} {
		res := callTool(t, cs, name, args)
		if !res.IsError || !strings.Contains(resultText(t, res), "must be a UUID") {
			t.Errorf("%s: result = %+v", name, res)
		}
	}
}

func TestWriteToolsUnauthorized(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, http.StatusUnauthorized, `{"msg":"no"}`) })
	cs := newTestClient(t, h, Options{})
	for name, args := range map[string]map[string]any{
		"create_memo":   {"title": "t"},
		"update_memo":   {"id": testMemoID, "title": "t"},
		"edit_memo":     {"id": testMemoID, "old_string": "a", "new_string": "b"},
		"delete_memo":   {"id": testMemoID},
		"set_memo_tags": {"id": testMemoID, "tags": []string{}},
		"publish_memo":  {"id": testMemoID, "state": "private"},
	} {
		res := callTool(t, cs, name, args)
		if !res.IsError || !strings.Contains(resultText(t, res), "webmemo login") {
			t.Errorf("%s: want tool error mentioning login, got %+v", name, res)
		}
	}
}
