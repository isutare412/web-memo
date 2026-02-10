package handlers

import (
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/web/gen"
)

// UploadURLToWeb converts a model.UploadURL to the generated CreateUploadURLResponse.
func UploadURLToWeb(u *model.UploadURL) gen.CreateUploadURLResponse {
	return gen.CreateUploadURLResponse{
		ImageID:       u.ImageID,
		UploadURL:     u.URL,
		UploadHeaders: u.Headers,
		ExpiresAt:     u.ExpiresAt,
	}
}

// ImageToWeb converts a model.Image to the generated ImageStatus.
func ImageToWeb(img *model.Image) gen.ImageStatus {
	return gen.ImageStatus{
		ID:         img.ID,
		State:      img.State,
		Original:   ImageDataToWeb(img.Original),
		Downscaled: ImageDataToWeb(img.Downscaled),
	}
}

// ImageDataToWeb converts a model.ImageData to the generated ImageVariant.
// Returns nil if the input is nil.
func ImageDataToWeb(d *model.ImageData) *gen.ImageVariant {
	if d == nil {
		return nil
	}
	return &gen.ImageVariant{
		URL:    d.URL,
		Format: d.Format,
	}
}
