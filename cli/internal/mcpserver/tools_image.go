package mcpserver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/isutare412/web-memo/cli/internal/apiclient/gen"
	"github.com/isutare412/web-memo/cli/internal/session"
)

// maxImageBytes caps the size of a file upload_image reads into memory.
const maxImageBytes = 50 << 20

// imageStateReady is the state of an image that finished processing.
const imageStateReady gen.ImageState = "READY"

func (t *tools) registerImage(srv *mcp.Server) {
	mcp.AddTool(srv, &mcp.Tool{
		Name: "upload_image",
		Description: "Upload a local image file (JPEG, PNG, WEBP, AVIF or HEIC, up to 50 MiB) and wait until it is processed. " +
			"path is a file path on the machine running `webmemo mcp`, not a URL. " +
			"Returns imageId, url (display version), originalUrl and markdown. " +
			"Uploading alone does not attach the image to any memo: insert the returned markdown into a memo " +
			"with edit_memo, create_memo or update_memo.",
		Annotations: additiveAnnotations,
	}, t.uploadImage)
}

type uploadImageInput struct {
	Path string `json:"path" jsonschema:"path of the image file on the machine running webmemo mcp"`
}

type uploadImageResult struct {
	ImageID     string `json:"imageId"`
	Markdown    string `json:"markdown"`
	URL         string `json:"url"`
	OriginalURL string `json:"originalUrl"`
}

func (t *tools) uploadImage(ctx context.Context, _ *mcp.CallToolRequest, in uploadImageInput) (*mcp.CallToolResult, any, error) {
	data, format, err := readImageFile(in.Path)
	if err != nil {
		return errorResult(err), nil, nil
	}

	created, err := t.session.API().CreateUploadURLWithResponse(ctx, gen.CreateUploadURLJSONRequestBody{
		FileName: filepath.Base(in.Path),
		Format:   format,
	})
	if err != nil {
		return errorResult(err), nil, nil
	}
	if err := session.CheckStatus(created.StatusCode(), created.Body); err != nil {
		return errorResult(err), nil, nil
	}
	if created.JSON200 == nil {
		return errorResult(errUnexpectedBody), nil, nil
	}
	target := created.JSON200

	if err := t.putToStorage(ctx, target.UploadURL, target.UploadHeaders, data); err != nil {
		return errorResult(err), nil, nil
	}

	wait := true
	status, err := t.session.API().GetImageStatusWithResponse(ctx, target.ImageID, &gen.GetImageStatusParams{WaitUntilProcessed: &wait})
	if err != nil {
		return errorResult(err), nil, nil
	}
	if err := session.CheckStatus(status.StatusCode(), status.Body); err != nil {
		return errorResult(err), nil, nil
	}
	if status.JSON200 == nil {
		return errorResult(errUnexpectedBody), nil, nil
	}
	img := status.JSON200
	if img.State != imageStateReady {
		return errorResult(fmt.Errorf("image %s processing did not finish: state is %s", target.ImageID, img.State)), nil, nil
	}

	// Same selection as the web UI (ui/src/lib/imageUpload.ts).
	var display, original string
	if img.Original != nil {
		original = img.Original.URL
	}
	display = original
	if img.Downscaled != nil {
		display = img.Downscaled.URL
	}
	if original == "" {
		original = display
	}
	return respond(uploadImageResult{
		ImageID:     target.ImageID,
		Markdown:    imageMarkdown(display, original),
		URL:         display,
		OriginalURL: original,
	})
}

// readImageFile reads a regular image file after checking its size and
// detecting its format from the content.
func readImageFile(path string) ([]byte, gen.ImageFormat, error) {
	if path == "" {
		return nil, "", errors.New("path must not be empty")
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, "", fmt.Errorf("cannot read image file: %w", err)
	}
	switch {
	case !info.Mode().IsRegular():
		return nil, "", fmt.Errorf("%s is not a regular file", path)
	case info.Size() == 0:
		return nil, "", fmt.Errorf("%s is empty", path)
	case info.Size() > maxImageBytes:
		return nil, "", fmt.Errorf("%s is too large (%d bytes): the limit is %d MiB", path, info.Size(), maxImageBytes>>20)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("cannot read image file: %w", err)
	}
	format, err := detectImageFormat(data[:min(len(data), imageHeadSize)])
	if err != nil {
		return nil, "", err
	}
	return data, format, nil
}

// putToStorage uploads data to the presigned URL. It uses the session's
// unauthenticated client so the Bearer token never reaches the storage host.
func (t *tools) putToStorage(ctx context.Context, url string, headers map[string]string, data []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("build upload request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := t.session.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("upload image: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("upload image: storage responded with status %d: %s", resp.StatusCode, msg)
	}
	return nil
}

// imageMarkdown builds the markdown for an uploaded image: a link to the
// original when it differs from the displayed (downscaled) version.
func imageMarkdown(displayURL, originalURL string) string {
	if displayURL == "" {
		displayURL = originalURL
	}
	if originalURL == "" {
		originalURL = displayURL
	}
	if displayURL != originalURL {
		return fmt.Sprintf("[![image](%s)](%s)", displayURL, originalURL)
	}
	return fmt.Sprintf("![image](%s)", displayURL)
}
