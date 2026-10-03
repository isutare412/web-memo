package mcpserver

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

const (
	testMemoID  = "11111111-1111-4111-8111-111111111111"
	testOwnerID = "22222222-2222-4222-8222-222222222222"
	testUserID  = "33333333-3333-4333-8333-333333333333"
)

func writeJSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func memosHandler(t *testing.T, got *url.Values, content string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos" {
			t.Errorf("path = %s", r.URL.Path)
		}
		*got = r.URL.Query()
		body, _ := json.Marshal(map[string]any{
			"page": 2, "pageSize": 5, "lastPage": 3, "totalMemoCount": 11,
			"memos": []map[string]any{{
				"id": testMemoID, "ownerId": testOwnerID, "title": "제목", "content": content,
				"tags": []string{"a", "b"}, "publishState": "private", "version": 1,
				"createTime": "2026-01-01T00:00:00Z", "updateTime": "2026-01-02T00:00:00Z",
				"scores": map[string]any{"bm25": 1.5, "rrf": 0.5, "semantic": 0.25},
			}},
		})
		writeJSON(w, http.StatusOK, string(body))
	})
}

type searchMemosResult struct {
	Page, PageSize, LastPage, TotalMemoCount int
	Memos                                    []struct {
		ID, Title, Preview, PublishState string
		Tags                             []string
		Scores                           map[string]float64
		Content                          *string
	}
}

func TestSearchMemosQueryMode(t *testing.T) {
	content := strings.Repeat("가", 300)
	var gotQuery url.Values
	cs := newTestClient(t, memosHandler(t, &gotQuery, content), Options{})

	res := callTool(t, cs, "search_memos", map[string]any{"query": "hello"})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	if want := (url.Values{"q": {"hello"}}); !reflect.DeepEqual(gotQuery, want) {
		t.Errorf("query = %v, want %v", gotQuery, want)
	}

	var out searchMemosResult
	if err := json.Unmarshal([]byte(resultText(t, res)), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out.Memos) != 1 {
		t.Fatalf("memos = %d", len(out.Memos))
	}
	m := out.Memos[0]
	if m.ID != testMemoID || m.Title != "제목" || m.PublishState != "private" || !reflect.DeepEqual(m.Tags, []string{"a", "b"}) {
		t.Errorf("memo = %+v", m)
	}
	if !utf8.ValidString(m.Preview) {
		t.Errorf("preview is not valid UTF-8")
	}
	if want := strings.Repeat("가", 200) + "…"; m.Preview != want {
		t.Errorf("preview = %d runes, want 200 runes + ellipsis", utf8.RuneCountInString(m.Preview))
	}
	if m.Content != nil {
		t.Errorf("search result must not include full content")
	}
	if m.Scores["rrf"] != 0.5 || m.Scores["bm25"] != 1.5 || m.Scores["semantic"] != 0.25 {
		t.Errorf("scores = %v", m.Scores)
	}
}

func TestSearchMemosListMode(t *testing.T) {
	var gotQuery url.Values
	cs := newTestClient(t, memosHandler(t, &gotQuery, "short"), Options{})

	res := callTool(t, cs, "search_memos", map[string]any{
		"tags": []string{"go", "mcp"}, "sort": "updateTime", "page": 2, "page_size": 5,
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	want := url.Values{"tag": {"go", "mcp"}, "sort": {"updateTime"}, "page": {"2"}, "pageSize": {"5"}}
	if !reflect.DeepEqual(gotQuery, want) {
		t.Errorf("query = %v, want %v", gotQuery, want)
	}

	var out searchMemosResult
	if err := json.Unmarshal([]byte(resultText(t, res)), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Page != 2 || out.PageSize != 5 || out.LastPage != 3 || out.TotalMemoCount != 11 {
		t.Errorf("paging = %+v", out)
	}
	if len(out.Memos) != 1 || out.Memos[0].Preview != "short" {
		t.Errorf("memos = %+v", out.Memos)
	}
}

func TestSearchMemosRejectsInvalidSort(t *testing.T) {
	cs := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("API must not be called")
	}), Options{})
	res := callTool(t, cs, "search_memos", map[string]any{"sort": "bogus"})
	if !res.IsError || !strings.Contains(resultText(t, res), "sort") {
		t.Errorf("want sort error, got %+v", res)
	}
}

func TestGetMemo(t *testing.T) {
	cs := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos/"+testMemoID {
			t.Errorf("path = %s", r.URL.Path)
		}
		writeJSON(w, http.StatusOK, `{"id":"`+testMemoID+`","ownerId":"`+testOwnerID+`","title":"T","content":"full body",
			"tags":[],"publishState":"private","version":3,"createTime":"2026-01-01T00:00:00Z","updateTime":"2026-01-02T00:00:00Z"}`)
	}), Options{})

	res := callTool(t, cs, "get_memo", map[string]any{"id": testMemoID})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	var memo struct {
		Content string
		Version int
	}
	if err := json.Unmarshal([]byte(resultText(t, res)), &memo); err != nil {
		t.Fatal(err)
	}
	if memo.Content != "full body" || memo.Version != 3 {
		t.Errorf("memo = %+v", memo)
	}
}

func TestGetMemoNotFound(t *testing.T) {
	cs := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{"msg":"memo not found"}`)
	}), Options{})

	res := callTool(t, cs, "get_memo", map[string]any{"id": testMemoID})
	if !res.IsError {
		t.Fatalf("IsError = false, want true")
	}
	if text := resultText(t, res); !strings.Contains(text, "memo not found") {
		t.Errorf("text = %q", text)
	}
}

func TestGetMemoInvalidID(t *testing.T) {
	cs := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("API must not be called")
	}), Options{})
	res := callTool(t, cs, "get_memo", map[string]any{"id": "not-a-uuid"})
	if !res.IsError || !strings.Contains(resultText(t, res), "not-a-uuid") {
		t.Errorf("want invalid id error, got %+v", res)
	}
}

func TestListMemoMembers(t *testing.T) {
	cs := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/memos/" + testMemoID + "/subscribers":
			writeJSON(w, http.StatusOK, `{"memoOwnerId":"`+testOwnerID+`","subscribers":[{"id":"`+testUserID+`","userName":"sub","photoUrl":"p","approved":true}]}`)
		case "/api/v1/memos/" + testMemoID + "/collaborators":
			writeJSON(w, http.StatusOK, `{"memoOwnerId":"`+testOwnerID+`","collaborators":[{"id":"`+testUserID+`","userName":"col","photoUrl":"p","isApproved":false}]}`)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}), Options{})

	res := callTool(t, cs, "list_memo_members", map[string]any{"id": testMemoID})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	var out struct {
		MemoOwnerID string
		Subscribers []struct {
			UserName string
			Approved bool
		}
		Collaborators []struct {
			UserName   string
			IsApproved bool
		}
	}
	if err := json.Unmarshal([]byte(resultText(t, res)), &out); err != nil {
		t.Fatal(err)
	}
	if out.MemoOwnerID != testOwnerID {
		t.Errorf("memoOwnerId = %q", out.MemoOwnerID)
	}
	if len(out.Subscribers) != 1 || out.Subscribers[0].UserName != "sub" || !out.Subscribers[0].Approved {
		t.Errorf("subscribers = %+v", out.Subscribers)
	}
	if len(out.Collaborators) != 1 || out.Collaborators[0].UserName != "col" || out.Collaborators[0].IsApproved {
		t.Errorf("collaborators = %+v", out.Collaborators)
	}
}

func TestSearchTags(t *testing.T) {
	cs := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tags" || r.URL.Query().Get("kw") != "go" {
			t.Errorf("request = %s", r.URL.String())
		}
		writeJSON(w, http.StatusOK, `["go","golang"]`)
	}), Options{})

	res := callTool(t, cs, "search_tags", map[string]any{"keyword": "go"})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	var tags []string
	if err := json.Unmarshal([]byte(resultText(t, res)), &tags); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(tags, []string{"go", "golang"}) {
		t.Errorf("tags = %v", tags)
	}
}

func TestPreview(t *testing.T) {
	tests := map[string]struct{ in, want string }{
		"empty":         {"", ""},
		"short":         {"abc", "abc"},
		"exactly limit": {strings.Repeat("a", 200), strings.Repeat("a", 200)},
		"over limit":    {strings.Repeat("a", 201), strings.Repeat("a", 200) + "…"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := preview(tc.in); got != tc.want {
				t.Errorf("preview() = %q, want %q", got, tc.want)
			}
		})
	}
}

// The generated client sends `tag=` for a nil tag list, which the API reads as
// a filter on the empty tag name.
func TestSearchMemosWithoutTagsSendsNoTagParam(t *testing.T) {
	var gotQuery url.Values
	cs := newTestClient(t, memosHandler(t, &gotQuery, "short"), Options{})

	res := callTool(t, cs, "search_memos", map[string]any{})
	if res.IsError {
		t.Fatalf("unexpected error: %s", resultText(t, res))
	}
	if _, ok := gotQuery["tag"]; ok || len(gotQuery) != 0 {
		t.Errorf("query = %v, want empty", gotQuery)
	}
}
