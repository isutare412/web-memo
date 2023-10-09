package model

import (
	"time"

	"github.com/google/uuid"
)

type GoogleTokenResponse struct {
	AccessToken    string
	AccessTokenTTL time.Duration
	IDToken        string
	Scope          string
	TokenType      string
}

type GoogleIDToken struct {
	IssuedAt   time.Time
	ExpiresAt  time.Time
	Subject    string
	Email      string
	Name       string
	FamilyName string
	GivenName  string
	PictureURL string
}

type AppIDToken struct {
	UserID     uuid.UUID
	UserType   UserType
	Email      string
	UserName   string
	FamilyName string
	GivenName  string
	PhotoURL   string
}
