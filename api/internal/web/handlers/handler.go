package handlers

import (
	"time"

	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/web/gen"
)

var _ gen.ServerInterface = (*Handler)(nil)

type Handler struct {
	authService      port.AuthService
	memoService      port.MemoService
	imageService     port.ImageService
	pingers          []port.Pinger
	cookieExpiration time.Duration
	embeddingEnabled bool
}

func NewHandler(
	pingers []port.Pinger,
	authService port.AuthService,
	memoService port.MemoService,
	imageService port.ImageService,
	cookieExpiration time.Duration,
	embeddingEnabled bool,
) *Handler {
	return &Handler{
		authService:      authService,
		memoService:      memoService,
		imageService:     imageService,
		pingers:          pingers,
		cookieExpiration: cookieExpiration,
		embeddingEnabled: embeddingEnabled,
	}
}
