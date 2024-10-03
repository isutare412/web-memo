package backup

import (
	"time"

	"github.com/isutare412/web-memo/backup/internal/core/model"
)

type Config struct {
	DatabaseName     string
	DatabaseHost     string
	DatabasePort     int
	DatabaseUser     string
	DatabasePassword string

	Retention time.Duration
	Prefix    string
}

func (c Config) ToDatabaseBackupRequest() model.DatabaseBackupRequest {
	return model.DatabaseBackupRequest{
		Host:         c.DatabaseHost,
		Port:         c.DatabasePort,
		DatabaseName: c.DatabaseName,
		User:         c.DatabaseUser,
		Password:     c.DatabasePassword,
	}
}
