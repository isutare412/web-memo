package port

import (
	"github.com/isutare412/web-memo/api/internal/core/model"
)

type JWTClient interface {
	ParseGoogleIDTokenUnverified(tokenString string) (*model.GoogleIDToken, error)
}
