package http

import (
	"time"

	"github.com/isutare412/imageer/pkg/images"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type createUploadURLRequest struct {
	FileName string        `json:"fileName" validate:"required"`
	Format   images.Format `json:"format" validate:"required"`
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

type getImageStatusResponse struct {
	ID       string            `json:"id"`
	State    images.State      `json:"state"`
	URL      string            `json:"url"`
	Format   images.Format     `json:"format"`
	Variants []imageVariantDTO `json:"variants"`
}

type imageVariantDTO struct {
	ID         string        `json:"id"`
	PresetName string        `json:"presetName"`
	URL        string        `json:"url"`
	Format     images.Format `json:"format"`
	State      string        `json:"state"`
}

func (r *getImageStatusResponse) fromModel(img *model.Image) {
	r.ID = img.ID
	r.State = img.State
	r.URL = img.URL
	r.Format = img.Format
	r.Variants = make([]imageVariantDTO, 0, len(img.Variants))
	for _, v := range img.Variants {
		r.Variants = append(r.Variants, imageVariantDTO{
			ID:         v.ID,
			PresetName: v.PresetName,
			URL:        v.URL,
			Format:     v.Format,
			State:      v.State,
		})
	}
}
