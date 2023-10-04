package auth

import "time"

type Config struct {
	Google            GoogleConfig
	OAuthStateTimeout time.Duration
}

type GoogleConfig struct {
	OAuthEndpoint     string
	OAuthClientID     string
	OAuthCallbackPath string
}
