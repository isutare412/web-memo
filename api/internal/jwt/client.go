package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) ParseGoogleIDTokenUnverified(tokenString string) (*model.GoogleIDToken, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &googleIDTokenClaim{})
	if err != nil {
		return nil, fmt.Errorf("parsing google ID token: %w", err)
	}

	claims, ok := token.Claims.(*googleIDTokenClaim)
	if !ok {
		return nil, fmt.Errorf("unexpected token claim type")
	}

	return claims.ToGoogleIDToken(), nil
}
