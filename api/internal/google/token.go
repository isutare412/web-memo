package google

import (
	"time"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type googleOAuthTokens struct {
	AccessToken           string `json:"access_token"`
	AccessTokenTTLSeconds int    `json:"expires_in"`
	IDToken               string `json:"id_token"`
	Scope                 string `json:"scope"`
	TokenType             string `json:"token_type"`
}

func (t *googleOAuthTokens) ToTokenResponse() model.GoogleTokenResponse {
	return model.GoogleTokenResponse{
		AccessToken:    t.AccessToken,
		AccessTokenTTL: time.Second * time.Duration(t.AccessTokenTTLSeconds),
		IDToken:        t.IDToken,
		Scope:          t.Scope,
		TokenType:      t.TokenType,
	}
}
