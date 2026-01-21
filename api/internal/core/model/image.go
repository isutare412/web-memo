package model

import (
	"time"

	"github.com/isutare412/imageer/pkg/images"
)

type UploadURL struct {
	ImageID   string
	URL       string
	Headers   map[string]string
	ExpiresAt time.Time
}

type ImageData struct {
	URL    string
	Format images.Format
}

type Image struct {
	ID         string
	State      images.State
	Original   *ImageData
	Downscaled *ImageData
}
