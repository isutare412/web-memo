package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

const webMemoIssuer = "web-memo"

type googleIDTokenClaims struct {
	jwt.RegisteredClaims
	Email         string `json:"email"`
	EmailVerified *bool  `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	PicturlURL    string `json:"picture"`
	ProfileURL    string `json:"profile"`
}

func (c *googleIDTokenClaims) ToGoogleIDToken() *model.GoogleIDToken {
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

type appClaims struct {
	jwt.RegisteredClaims
	UserID     string `json:"user_id"`
	UserType   string `json:"user_type"`
	Email      string `json:"email"`
	UserName   string `json:"name"`
	FamilyName string `json:"family_name,omitempty"`
	GivenName  string `json:"given_name,omitempty"`
	PhotoURL   string `json:"photo_url,omitempty"`
}

func newAppClaims(token *model.AppIDToken, ttl time.Duration) appClaims {
	now := time.Now()

	return appClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    webMemoIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
		UserID:     token.UserID.String(),
		UserType:   string(token.UserType),
		Email:      token.Email,
		UserName:   token.UserName,
		FamilyName: token.FamilyName,
		GivenName:  token.GivenName,
		PhotoURL:   token.PhotoURL,
	}
}

func (c *appClaims) toAppIDToken() (*model.AppIDToken, error) {
	userID, err := uuid.Parse(c.UserID)
	if err != nil {
		return nil, fmt.Errorf("parsing user ID: %w", err)
	}

	return &model.AppIDToken{
		UserID:     userID,
		UserType:   model.UserType(c.UserType),
		Email:      c.Email,
		UserName:   c.UserName,
		FamilyName: c.FamilyName,
		GivenName:  c.GivenName,
		PhotoURL:   c.PhotoURL,
	}, nil
}
