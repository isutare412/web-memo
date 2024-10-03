package config

import (
	"time"

	"github.com/isutare412/web-memo/backup/internal/aws"
	"github.com/isutare412/web-memo/backup/internal/core/service/backup"
)

type Config struct {
	ProcessTimeout time.Duration  `koanf:"process-timeout" validate:"required"`
	AWS            AWSConfig      `koanf:"aws"`
	Backup         BackupConfig   `koanf:"backup"`
	Postgres       PostgresConfig `koanf:"postgres"`
}

func (c *Config) ToAWSS3Config() aws.S3Config {
	return aws.S3Config{
		AccessKey: c.AWS.AccessKey,
		Secret:    c.AWS.Secret,
		Bucket:    c.AWS.S3.BackupBucket,
	}
}

func (c *Config) ToBackupConfig() backup.Config {
	return backup.Config{
		DatabaseName:     c.Postgres.Database,
		DatabaseHost:     c.Postgres.Host,
		DatabasePort:     c.Postgres.Port,
		DatabaseUser:     c.Postgres.User,
		DatabasePassword: c.Postgres.Password,
		Retention:        c.Backup.Retention,
		Prefix:           c.Backup.Prefix,
	}
}

type AWSConfig struct {
	AccessKey string   `koanf:"access-key" validate:"required"`
	Secret    string   `koanf:"secret" validate:"required"`
	S3        S3Config `koanf:"s3"`
}

type S3Config struct {
	BackupBucket string `koanf:"backup-bucket" validate:"required"`
}

type BackupConfig struct {
	Prefix    string        `koanf:"prefix" validate:"endswith=/"`
	Retention time.Duration `koanf:"retention" validate:"required"`
}

type PostgresConfig struct {
	Host     string `koanf:"host" validate:"required"`
	Port     int    `koanf:"port" validate:"required"`
	User     string `koanf:"user" validate:"required"`
	Password string `koanf:"password" validate:"required"`
	Database string `koanf:"database" validate:"required"`
}
