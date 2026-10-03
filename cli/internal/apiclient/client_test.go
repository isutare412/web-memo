package apiclient_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isutare412/web-memo/cli/internal/apiclient/gen"
)

const listMemosBody = `{
  "memos": [
    {
      "id": "0b5f6b0e-3c1e-4a53-8a52-6a1a7f7a3c11",
      "ownerId": "5b0c3f1e-8d0e-4d8e-9d3b-2f6f1f0a9c22",
      "title": "k8s notes",
      "content": "# kubernetes",
      "tags": ["k8s"],
      "version": 3,
      "publishState": "shared",
      "createTime": "2026-10-01T09:00:00Z",
      "updateTime": "2026-10-02T09:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 10,
  "lastPage": 1,
  "totalMemoCount": 1
}`

func TestListMemosDecodes(t *testing.T) {
	var gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(listMemosBody))
	}))
	defer srv.Close()

	client, err := gen.NewClientWithResponses(srv.URL, gen.WithHTTPClient(srv.Client()))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	q := "k8s"
	resp, err := client.ListMemosWithResponse(context.Background(), &gen.ListMemosParams{SearchQuery: &q})
	if err != nil {
		t.Fatalf("list memos: %v", err)
	}

	if gotPath != "/api/v1/memos" {
		t.Errorf("path = %q, want /api/v1/memos", gotPath)
	}
	if gotQuery != "k8s" {
		t.Errorf("q = %q, want k8s", gotQuery)
	}
	if resp.JSON200 == nil {
		t.Fatalf("JSON200 is nil; status=%s body=%s", resp.Status(), resp.Body)
	}
	if len(resp.JSON200.Memos) != 1 {
		t.Fatalf("memos = %d, want 1", len(resp.JSON200.Memos))
	}
	if got := resp.JSON200.Memos[0].PublishState; got != gen.PublishStateShared {
		t.Errorf("publishState = %q, want %q", got, gen.PublishStateShared)
	}
}
