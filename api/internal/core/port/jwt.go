package port

import (
	"github.com/isutare412/web-memo/api/internal/core/model"
)

type JWTClient interface {
	ParseGoogleIDTokenUnverified(tokenString string) (*model.GoogleIDToken, error)
	SignAppIDToken(*model.AppIDToken) (token *model.AppIDToken, tokenString string, err error)
	VerifyAppIDToken(tokenString string) (*model.AppIDToken, error)
}
