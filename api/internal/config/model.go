package config

import (
	"time"

	"github.com/isutare412/web-memo/api/internal/log"
	"github.com/isutare412/web-memo/api/internal/postgres"
	"github.com/isutare412/web-memo/api/internal/redis"
)

type Config struct {
	Wire     WireConfig     `mapstructure:"wire"`
	Log      LogConfig      `mapstructure:"log"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

func (c *Config) ToLogConfig() log.Config {
	return log.Config(c.Log)
}

func (c *Config) ToPostgresConfig() postgres.Config {
	return postgres.Config(c.Postgres)
}

func (c *Config) ToRedisConfig() redis.Config {
	return redis.Config(c.Redis)
}

type WireConfig struct {
	InitializeTimeout time.Duration `mapstructure:"initializeTimeout" validate:"required"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdownTimeout" validate:"required"`
}

type LogConfig struct {
	Format log.Format `mapstructure:"format" validate:"required,oneof=json text"`
	Level  log.Level  `mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Caller bool       `mapstructure:"caller"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Database string `mapstructure:"database" validate:"required"`
	QueryLog bool   `mapstructure:"queryLog"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
}
