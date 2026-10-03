package mcpserver

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

const testImageID = "img-1"

var pngBytes = append([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1A, '\n'}, []byte("fake png payload")...)

func writeTempFile(t *testing.T, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// imageAPI is a fake web-memo API plus a fake presigned S3 endpoint at /s3.
type imageAPI struct {
	t           *testing.T
	statusBody  string
	s3Status    int
	mu          sync.Mutex
	s3Body      []byte
	s3Header    http.Header
	s3Method    string
	uploadReqs  []map[string]any
	statusQuery string
	apiCalls    int
	s3Calls     int
}

func (f *imageAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	switch {
	case r.URL.Path == "/s3":
		f.s3Calls++
		f.s3Method = r.Method
		f.s3Header = r.Header.Clone()
		f.s3Body, _ = io.ReadAll(r.Body)
		w.WriteHeader(f.s3Status)
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/images/upload-url":
		f.apiCalls++
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		f.uploadReqs = append(f.uploadReqs, body)
		writeJSON(w, http.StatusOK, `{"imageId":"`+testImageID+`","uploadUrl":"http://`+r.Host+`/s3","uploadHeaders":{"Content-Type":"image/png","x-amz-meta-a":"b"},"expiresAt":"2030-01-01T00:00:00Z"}`)
	case r.Method == http.MethodGet && r.URL.Path == "/api/v1/images/"+testImageID+"/status":
		f.apiCalls++
		f.statusQuery = r.URL.RawQuery
		writeJSON(w, http.StatusOK, f.statusBody)
	default:
		f.t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected", http.StatusInternalServerError)
	}
}

func newImageAPI(t *testing.T, statusBody string) *imageAPI {
	return &imageAPI{t: t, statusBody: statusBody, s3Status: http.StatusOK}
}

const readyBody = `{"id":"` + testImageID + `","state":"READY","original":{"url":"https://cdn/o.png","format":"PNG"},"downscaled":{"url":"https://cdn/d.webp","format":"WEBP"}}`

func TestUploadImageFlow(t *testing.T) {
	api := newImageAPI(t, readyBody)
	cs := newTestClient(t, api, Options{})
	path := writeTempFile(t, "shot.png", pngBytes)

	res := callTool(t, cs, "upload_image", map[string]any{"path": path})
	if res.IsError {
		t.Fatalf("unexpected tool error: %s", resultText(t, res))
	}

	var out map[string]string
	if err := json.Unmarshal([]byte(resultText(t, res)), &out); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"imageId":     testImageID,
		"markdown":    "[![image](https://cdn/d.webp)](https://cdn/o.png)",
		"url":         "https://cdn/d.webp",
		"originalUrl": "https://cdn/o.png",
	}
	for k, v := range want {
		if out[k] != v {
			t.Errorf("%s = %q, want %q", k, out[k], v)
		}
	}

	if len(api.uploadReqs) != 1 || api.uploadReqs[0]["fileName"] != "shot.png" || api.uploadReqs[0]["format"] != "PNG" {
		t.Errorf("createUploadURL requests = %v", api.uploadReqs)
	}
	if api.s3Method != http.MethodPut {
		t.Errorf("s3 method = %q, want PUT", api.s3Method)
	}
	if string(api.s3Body) != string(pngBytes) {
		t.Errorf("s3 body = %q, want file bytes", api.s3Body)
	}
	if got := api.s3Header.Get("Content-Type"); got != "image/png" {
		t.Errorf("s3 Content-Type = %q, want image/png", got)
	}
	if got := api.s3Header.Get("X-Amz-Meta-A"); got != "b" {
		t.Errorf("s3 x-amz-meta-a = %q, want b", got)
	}
	if got := api.s3Header.Get("Authorization"); got != "" {
		t.Errorf("s3 request leaked Authorization header %q", got)
	}
	if api.statusQuery != "waitUntilProcessed=true" {
		t.Errorf("status query = %q, want waitUntilProcessed=true", api.statusQuery)
	}
}

func TestUploadImageWithoutDownscaled(t *testing.T) {
	api := newImageAPI(t, `{"id":"`+testImageID+`","state":"READY","original":{"url":"https://cdn/o.png","format":"PNG"}}`)
	cs := newTestClient(t, api, Options{})
	res := callTool(t, cs, "upload_image", map[string]any{"path": writeTempFile(t, "a.png", pngBytes)})
	if res.IsError {
		t.Fatalf("unexpected tool error: %s", resultText(t, res))
	}
	if text := resultText(t, res); !strings.Contains(text, `"markdown": "![image](https://cdn/o.png)"`) {
		t.Errorf("result = %s, want plain image markdown", text)
	}
}

func TestUploadImageMissingFile(t *testing.T) {
	api := newImageAPI(t, readyBody)
	cs := newTestClient(t, api, Options{})
	res := callTool(t, cs, "upload_image", map[string]any{"path": filepath.Join(t.TempDir(), "nope.png")})
	if !res.IsError {
		t.Fatalf("want tool error, got %s", resultText(t, res))
	}
	if api.apiCalls != 0 || api.s3Calls != 0 {
		t.Errorf("no request expected, got api=%d s3=%d", api.apiCalls, api.s3Calls)
	}
}

func TestUploadImageRejectsBadInput(t *testing.T) {
	dir := t.TempDir()
	big := filepath.Join(dir, "big.png")
	f, err := os.Create(big)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(maxImageBytes + 1); err != nil { // sparse file
		t.Fatal(err)
	}
	_ = f.Close()

	tests := []struct {
		name, path, wantMsg string
	}{
		{"directory", dir, "not a regular file"},
		{"text file", writeTempFile(t, "notes.txt", []byte("hello world, definitely not an image")), "unsupported image format"},
		{"empty file", writeTempFile(t, "empty.png", nil), "empty"},
		{"too large", big, "too large"},
		{"empty path", "", "path"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newImageAPI(t, readyBody)
			cs := newTestClient(t, api, Options{})
			res := callTool(t, cs, "upload_image", map[string]any{"path": tt.path})
			if !res.IsError {
				t.Fatalf("want tool error, got %s", resultText(t, res))
			}
			if text := resultText(t, res); !strings.Contains(text, tt.wantMsg) {
				t.Errorf("error = %q, want it to contain %q", text, tt.wantMsg)
			}
			if api.apiCalls != 0 || api.s3Calls != 0 {
				t.Errorf("no request expected, got api=%d s3=%d", api.apiCalls, api.s3Calls)
			}
		})
	}
}

func TestUploadImageFailedState(t *testing.T) {
	api := newImageAPI(t, `{"id":"`+testImageID+`","state":"FAILED"}`)
	cs := newTestClient(t, api, Options{})
	res := callTool(t, cs, "upload_image", map[string]any{"path": writeTempFile(t, "a.png", pngBytes)})
	if !res.IsError {
		t.Fatalf("want tool error, got %s", resultText(t, res))
	}
	if text := resultText(t, res); !strings.Contains(text, "FAILED") {
		t.Errorf("error = %q, want it to mention the FAILED state", text)
	}
}

func TestUploadImageS3Rejects(t *testing.T) {
	api := newImageAPI(t, readyBody)
	api.s3Status = http.StatusForbidden
	cs := newTestClient(t, api, Options{})
	res := callTool(t, cs, "upload_image", map[string]any{"path": writeTempFile(t, "a.png", pngBytes)})
	if !res.IsError {
		t.Fatalf("want tool error, got %s", resultText(t, res))
	}
	if text := resultText(t, res); !strings.Contains(text, "403") {
		t.Errorf("error = %q, want it to mention status 403", text)
	}
	if api.statusQuery != "" {
		t.Errorf("image status must not be polled after a failed upload")
	}
}

func TestUploadImageAPIError(t *testing.T) {
	cs := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{"msg":"unsupported format"}`)
	}), Options{})
	res := callTool(t, cs, "upload_image", map[string]any{"path": writeTempFile(t, "a.png", pngBytes)})
	if !res.IsError || !strings.Contains(resultText(t, res), "unsupported format") {
		t.Errorf("want tool error carrying the API message, got %+v", res)
	}
}

func TestImageMarkdown(t *testing.T) {
	tests := []struct {
		name, display, original, want string
	}{
		{"different urls link to original", "https://x/d", "https://x/o", "[![image](https://x/d)](https://x/o)"},
		{"same url", "https://x/o", "https://x/o", "![image](https://x/o)"},
		{"no downscaled variant", "", "https://x/o", "![image](https://x/o)"},
		{"no original variant", "https://x/d", "", "![image](https://x/d)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := imageMarkdown(tt.display, tt.original); got != tt.want {
				t.Errorf("imageMarkdown(%q, %q) = %q, want %q", tt.display, tt.original, got, tt.want)
			}
		})
	}
}
