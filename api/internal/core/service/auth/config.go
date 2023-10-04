package auth

import "time"

type Config struct {
	Google            GoogleConfig  `mapstructure:"google"`
	OAuthStateTimeout time.Duration `mapstructure:"oauth-state-timeout" validate:"required"`
}

type GoogleConfig struct {
	OAuthEndpoint     string `mapstructure:"oauth-endpoint" validate:"required"`
	OAuthCallbackPath string `mapstructure:"oauth-callback-path" validate:"required"`
	OAuthClientID     string `mapstructure:"oauth-client-id" validate:"required"`
	OAuthClientSecret string `mapstructure:"oauth-client-secret" validate:"required"`
}
