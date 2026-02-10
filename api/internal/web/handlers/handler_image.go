package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/tracing"
	"github.com/isutare412/web-memo/api/internal/web/gen"
)

// CreateUploadURL creates a presigned URL for uploading an image.
func (h *Handler) CreateUploadURL(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.CreateUploadURL")
	defer span.End()

	var req gen.CreateUploadURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}

	uploadURL, err := h.imageService.CreateUploadURL(ctx, req.FileName, req.Format)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("creating upload URL: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, UploadURLToWeb(uploadURL))
}

// GetImageStatus returns the processing status of an image.
func (h *Handler) GetImageStatus(w http.ResponseWriter, r *http.Request, imageID gen.ImageIDPath, params gen.GetImageStatusParams) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.GetImageStatus")
	defer span.End()

	waitUntilProcessed := params.WaitUntilProcessed != nil && *params.WaitUntilProcessed

	image, err := h.imageService.GetImage(ctx, imageID, waitUntilProcessed)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("getting image: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, ImageToWeb(image))
}
