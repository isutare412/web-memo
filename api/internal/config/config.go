package config

import (
	"time"

	"github.com/isutare412/web-memo/api/internal/core/service/auth"
	"github.com/isutare412/web-memo/api/internal/cron"
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
	HTTP     HTTPConfig      `mapstructure:"http"`
	Postgres postgres.Config `mapstructure:"postgres"`
	Redis    redis.Config    `mapstructure:"redis"`
	Google   GoogleConfig    `mapstructure:"google"`
	JWT      jwt.Config      `mapstructure:"jwt"`
	OAuth    OAuthConfig     `mapstructure:"oauth"`
	Cron     CronConfig      `mapstructure:"cron"`
}

func (c *Config) ToLogConfig() log.Config {
	return c.Log
}

func (c *Config) ToHTTPConfig() http.Config {
	return http.Config{
		Port:                  c.HTTP.Port,
		CookieTokenExpiration: c.JWT.Expiration,
	}
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

func (c *Config) ToCronConfig() cron.Config {
	return cron.Config{
		TagCleanupInterval: c.Cron.TagCleanupInterval,
	}
}

func (c *Config) ToAuthServiceConfig() auth.Config {
	return auth.Config{
		Google: auth.GoogleConfig{
			OAuthEndpoint:     c.Google.Endpoints.OAuth,
			OAuthClientID:     c.Google.OAuth.ClientID,
			OAuthCallbackPath: c.OAuth.CallbackPath,
		},
		OAuthStateTimeout: c.OAuth.StateTimeout,
	}
}

type WireConfig struct {
	InitializeTimeout time.Duration `mapstructure:"initialize-timeout" validate:"required"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown-timeout" validate:"required"`
}

type HTTPConfig struct {
	Port int `mapstructure:"port" validate:"required"`
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

type OAuthConfig struct {
	StateTimeout time.Duration `mapstructure:"state-timeout" validate:"required,min=1m"`
	CallbackPath string        `mapstructure:"callback-path" validate:"required"`
}

type CronConfig struct {
	TagCleanupInterval time.Duration `mapstructure:"tag-cleanup-interval" validate:"required,min=1m"`
}
