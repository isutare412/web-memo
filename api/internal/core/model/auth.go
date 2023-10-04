package model

import "time"

type GoogleTokenResponse struct {
	AccessToken    string
	AccessTokenTTL time.Duration
	IDToken        string
	Scope          string
	TokenType      string
}
