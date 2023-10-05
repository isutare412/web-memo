package config

import (
	"time"

	"github.com/isutare412/web-memo/api/internal/core/service/auth"
	"github.com/isutare412/web-memo/api/internal/google"
	"github.com/isutare412/web-memo/api/internal/http"
	"github.com/isutare412/web-memo/api/internal/jwt"
	"github.com/isutare412/web-memo/api/internal/log"
	"github.com/isutare412/web-memo/api/internal/postgres"
	"github.com/isutare412/web-memo/api/internal/redis"
)

type Config struct {
	Wire     WireConfig      `mapstructure:"wire"`
	Log      log.Config      `mapstructure:"log"`
	HTTP     http.Config     `mapstructure:"http"`
	Postgres postgres.Config `mapstructure:"postgres"`
	Redis    redis.Config    `mapstructure:"redis"`
	Google   GoogleConfig    `mapstructure:"google"`
	JWT      jwt.Config      `mapstructure:"jwt"`
	Service  ServiceConfig   `mapstructure:"service"`
}

func (c *Config) ToLogConfig() log.Config {
	return c.Log
}

func (c *Config) ToHTTPConfig() http.Config {
	return c.HTTP
}

func (c *Config) ToPostgresConfig() postgres.Config {
	return c.Postgres
}

func (c *Config) ToRedisConfig() redis.Config {
	return c.Redis
}

func (c *Config) ToGoogleClientConfig() google.ClientConfig {
	return google.ClientConfig{
		TokenEndpoint:     c.Google.Endpoints.Token,
		OAuthClientID:     c.Google.OAuth.ClientID,
		OAuthClientSecret: c.Google.OAuth.ClientSecret,
	}
}

func (c *Config) ToJWTConfig() jwt.Config {
	return jwt.Config(c.JWT)
}

func (c *Config) ToAuthServiceConfig() auth.Config {
	return auth.Config{
		Google: auth.GoogleConfig{
			OAuthEndpoint:     c.Google.Endpoints.OAuth,
			OAuthClientID:     c.Google.OAuth.ClientID,
			OAuthCallbackPath: c.Service.Auth.GoogleCallbackPath,
		},
		OAuthStateTimeout: c.Service.Auth.OAuthStateTimeout,
	}
}

type WireConfig struct {
	InitializeTimeout time.Duration `mapstructure:"initialize-timeout" validate:"required"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown-timeout" validate:"required"`
}

type GoogleConfig struct {
	Endpoints GoogleEndpointsConfig `mapstructure:"endpoints"`
	OAuth     GoogleOAuthConfig     `mapstructure:"oauth"`
}

type GoogleEndpointsConfig struct {
	Token string `mapstructure:"token" validate:"required"`
	OAuth string `mapstructure:"oauth" validate:"required"`
}

type GoogleOAuthConfig struct {
	ClientID     string `mapstructure:"client-id" validate:"required"`
	ClientSecret string `mapstructure:"client-secret" validate:"required"`
}

type ServiceConfig struct {
	Auth AuthServiceConfig `mapstructure:"auth"`
}

type AuthServiceConfig struct {
	OAuthStateTimeout  time.Duration `mapstructure:"oauth-state-timeout" validate:"required"`
	GoogleCallbackPath string        `mapstructure:"google-callback-path" validate:"required"`
}
