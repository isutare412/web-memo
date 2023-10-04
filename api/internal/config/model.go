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
	Wire     WireConfig     `mapstructure:"wire"`
	Log      LogConfig      `mapstructure:"log"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Service  ServiceConfig  `mapstructure:"service"`
}

func (c *Config) ToLogConfig() log.Config {
	return log.Config(c.Log)
}

func (c *Config) ToHTTPConfig() http.Config {
	return http.Config(c.HTTP)
}

func (c *Config) ToPostgresConfig() postgres.Config {
	return postgres.Config(c.Postgres)
}

func (c *Config) ToRedisConfig() redis.Config {
	return redis.Config(c.Redis)
}

func (c *Config) ToAuthServiceConfig() auth.Config {
	return c.Service.Auth
}

type WireConfig struct {
	InitializeTimeout time.Duration `mapstructure:"initialize-timeout" validate:"required"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown-timeout" validate:"required"`
}

type LogConfig struct {
	Format log.Format `mapstructure:"format" validate:"required,oneof=json text"`
	Level  log.Level  `mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Caller bool       `mapstructure:"caller"`
}

type HTTPConfig struct {
	Port int `json:"port" validate:"required"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Database string `mapstructure:"database" validate:"required"`
	QueryLog bool   `mapstructure:"query-log"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
}

type ServiceConfig struct {
	Auth auth.Config `mapstructure:"auth"`
}
