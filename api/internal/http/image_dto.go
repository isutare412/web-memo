package http

import (
	"time"

	"github.com/isutare412/imageer/pkg/images"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type createUploadURLRequest struct {
	FileName string        `json:"fileName"`
	Format   images.Format `json:"format"`
}

type createUploadURLResponse struct {
	ImageID       string            `json:"imageId"`
	UploadURL     string            `json:"uploadUrl"`
	UploadHeaders map[string]string `json:"uploadHeaders"`
	ExpiresAt     time.Time         `json:"expiresAt"`
}

func (r *createUploadURLResponse) fromModel(u *model.UploadURL) {
	r.ImageID = u.ImageID
	r.UploadURL = u.URL
	r.UploadHeaders = u.Headers
	r.ExpiresAt = u.ExpiresAt
}

type imageDataDTO struct {
	URL    string        `json:"url"`
	Format images.Format `json:"format"`
}

type getImageStatusResponse struct {
	ID         string        `json:"id"`
	State      images.State  `json:"state"`
	Original   *imageDataDTO `json:"original"`
	Downscaled *imageDataDTO `json:"downscaled"`
}

func (r *getImageStatusResponse) fromModel(img *model.Image) {
	r.ID = img.ID
	r.State = img.State
	if img.Original != nil {
		r.Original = &imageDataDTO{
			URL:    img.Original.URL,
			Format: img.Original.Format,
		}
	}
	if img.Downscaled != nil {
		r.Downscaled = &imageDataDTO{
			URL:    img.Downscaled.URL,
			Format: img.Downscaled.Format,
		}
	}
}
