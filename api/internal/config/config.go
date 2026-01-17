package config

import (
	"time"

	"github.com/isutare412/web-memo/api/internal/core/service/auth"
	"github.com/isutare412/web-memo/api/internal/core/service/image"
	"github.com/isutare412/web-memo/api/internal/cron"
	"github.com/isutare412/web-memo/api/internal/google"
	"github.com/isutare412/web-memo/api/internal/http"
	"github.com/isutare412/web-memo/api/internal/imageer"
	"github.com/isutare412/web-memo/api/internal/jwt"
	"github.com/isutare412/web-memo/api/internal/log"
	"github.com/isutare412/web-memo/api/internal/postgres"
	"github.com/isutare412/web-memo/api/internal/redis"
)

type Config struct {
	Wire     WireConfig      `koanf:"wire"`
	Log      log.Config      `koanf:"log"`
	HTTP     HTTPConfig      `koanf:"http"`
	Postgres postgres.Config `koanf:"postgres"`
	Redis    redis.Config    `koanf:"redis"`
	Google   GoogleConfig    `koanf:"google"`
	JWT      jwt.Config      `koanf:"jwt"`
	OAuth    OAuthConfig     `koanf:"oauth"`
	Cron     CronConfig      `koanf:"cron"`
	Imageer  ImageerConfig   `koanf:"imageer"`
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

func (c *Config) ToImageerConfig() imageer.Config {
	return imageer.Config{
		BaseURL: c.Imageer.BaseURL,
		APIKey:  c.Imageer.APIKey,
	}
}

func (c *Config) ToImageServiceConfig() image.Config {
	return image.Config{
		ProjectID: c.Imageer.ProjectID,
	}
}

type WireConfig struct {
	InitializeTimeout time.Duration `koanf:"initialize-timeout" validate:"required"`
	ShutdownTimeout   time.Duration `koanf:"shutdown-timeout" validate:"required"`
}

type HTTPConfig struct {
	Port int `koanf:"port" validate:"required"`
}

type GoogleConfig struct {
	Endpoints GoogleEndpointsConfig `koanf:"endpoints"`
	OAuth     GoogleOAuthConfig     `koanf:"oauth"`
}

type GoogleEndpointsConfig struct {
	Token string `koanf:"token" validate:"required"`
	OAuth string `koanf:"oauth" validate:"required"`
}

type GoogleOAuthConfig struct {
	ClientID     string `koanf:"client-id" validate:"required"`
	ClientSecret string `koanf:"client-secret" validate:"required"`
}

type OAuthConfig struct {
	StateTimeout time.Duration `koanf:"state-timeout" validate:"required,min=1m"`
	CallbackPath string        `koanf:"callback-path" validate:"required"`
}

type CronConfig struct {
	TagCleanupInterval time.Duration `koanf:"tag-cleanup-interval" validate:"required,min=1m"`
}

type ImageerConfig struct {
	BaseURL   string `koanf:"base-url" validate:"required,url"`
	APIKey    string `koanf:"api-key" validate:"required"`
	ProjectID string `koanf:"project-id" validate:"required"`
}
