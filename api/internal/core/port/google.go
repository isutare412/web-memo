package port

import (
	"context"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type GoogleClient interface {
	ExchangeAuthCode(ctx context.Context, code, redirectURI string) (model.GoogleTokenResponse, error)
}
