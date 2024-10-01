package port

import (
	"context"

	"github.com/isutare412/web-memo/backup/internal/core/model"
)

type BackupExecutor interface {
	BackupDatabase(context.Context, model.DatabaseBackupRequest) error
}
