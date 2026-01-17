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

type ImageVariant struct {
	ID         string
	PresetName string
	URL        string
	Format     images.Format
	State      string
}

type Image struct {
	ID       string
	State    images.State
	URL      string
	Format   images.Format
	Variants []ImageVariant
}
