package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type googleIDTokenClaim struct {
	jwt.RegisteredClaims
	Email         string `json:"email"`
	EmailVerified *bool  `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	PicturlURL    string `json:"picture"`
	ProfileURL    string `json:"profile"`
}

func (c *googleIDTokenClaim) ToGoogleIDToken() *model.GoogleIDToken {
	return &model.GoogleIDToken{
		IssuedAt:   c.IssuedAt.Time,
		ExpiresAt:  c.ExpiresAt.Time,
		Subject:    c.Subject,
		Email:      c.Email,
		Name:       c.Name,
		FamilyName: c.FamilyName,
		GivenName:  c.GivenName,
		PictureURL: c.PicturlURL,
	}
}
