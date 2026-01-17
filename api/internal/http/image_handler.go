package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/validate"
)

type imageHandler struct {
	imageService port.ImageService
}

func newImageHandler(imageService port.ImageService) *imageHandler {
	return &imageHandler{
		imageService: imageService,
	}
}

func (h *imageHandler) router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/upload-url", h.createUploadURL)
	r.Get("/{imageID}/status", h.getImageStatus)

	return r
}

func (h *imageHandler) createUploadURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createUploadURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}
	if err := validate.Struct(&req); err != nil {
		responseError(w, r, fmt.Errorf("validating request body: %w", err))
		return
	}

	uploadURL, err := h.imageService.CreateUploadURL(ctx, req.FileName, req.Format)
	if err != nil {
		responseError(w, r, fmt.Errorf("creating upload URL: %w", err))
		return
	}

	var resp createUploadURLResponse
	resp.fromModel(uploadURL)
	responseJSON(w, &resp)
}

func (h *imageHandler) getImageStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	imageIDStr := chi.URLParam(r, "imageID")
	if imageIDStr == "" {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "need imageID",
		})
		return
	}

	if _, err := uuid.Parse(imageIDStr); err != nil {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "invalid imageID format",
		})
		return
	}

	waitUntilProcessed := r.URL.Query().Get("waitUntilProcessed") == "true"

	image, err := h.imageService.GetImage(ctx, imageIDStr, waitUntilProcessed)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting image: %w", err))
		return
	}

	var resp getImageStatusResponse
	resp.fromModel(image)
	responseJSON(w, &resp)
}
