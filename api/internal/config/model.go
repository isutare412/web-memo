package config

import (
	"time"

	"github.com/isutare412/web-memo/api/internal/core/service/auth"
	"github.com/isutare412/web-memo/api/internal/http"
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

func (c *Config) ToAuthServiceConfig() auth.Config {
	return c.Service.Auth
}

type WireConfig struct {
	InitializeTimeout time.Duration `mapstructure:"initialize-timeout" validate:"required"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown-timeout" validate:"required"`
}

type ServiceConfig struct {
	Auth auth.Config `mapstructure:"auth"`
}
