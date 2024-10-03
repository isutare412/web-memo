package wire

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/isutare412/web-memo/backup/internal/aws"
	"github.com/isutare412/web-memo/backup/internal/backup"
	"github.com/isutare412/web-memo/backup/internal/config"
	corebackup "github.com/isutare412/web-memo/backup/internal/core/service/backup"
)

type App struct {
	cfg *config.Config

	backupService *corebackup.Service
}

func NewApp(cfg *config.Config) (*App, error) {
	backupExecutor := backup.NewExecutor()

	s3Client, err := aws.NewS3Client(cfg.ToAWSS3Config())
	if err != nil {
		return nil, fmt.Errorf("creating aws s3 client: %w", err)
	}

	backupService := corebackup.NewService(cfg.ToBackupConfig(), backupExecutor, s3Client)

	return &App{
		cfg:           cfg,
		backupService: backupService,
	}, nil
}

func (a *App) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.ProcessTimeout)
	defer cancel()

	if err := a.backupService.BackupDatabase(ctx); err != nil {
		slog.Error("failed to backup database", "error", err)
	} else {
		slog.Info("backup database complete")
	}

	if err := a.backupService.PruneDatabaseBackups(ctx); err != nil {
		slog.Error("failed to prune database backups", "error", err)
	} else {
		slog.Info("prune data backup complete")
	}
}
