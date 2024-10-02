package port

import (
	"context"
	"io"

	"github.com/isutare412/web-memo/backup/internal/core/model"
)

type BackupExecutor interface {
	BackupDatabase(context.Context, model.DatabaseBackupRequest) error
	ReadFile(ctx context.Context, fileName string) (io.ReadCloser, error)
}
