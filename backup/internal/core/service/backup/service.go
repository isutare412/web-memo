package backup

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"slices"
	"time"

	"github.com/isutare412/web-memo/backup/internal/core/port"
)

var regexPatternDateTime = regexp.MustCompile(`\d{14}`)

type Service struct {
	backupExecutor port.BackupExecutor
	objectStorage  port.ObjectStorage

	cfg Config
}

func NewService(
	cfg Config,
	backupExecutor port.BackupExecutor,
	objectStorage port.ObjectStorage,
) *Service {
	return &Service{
		backupExecutor: backupExecutor,
		objectStorage:  objectStorage,
		cfg:            cfg,
	}
}

func (s *Service) BackupDatabase(ctx context.Context) error {
	backup, err := s.backupExecutor.BackupDatabase(ctx, s.cfg.ToDatabaseBackupRequest())
	if err != nil {
		return fmt.Errorf("backup database: %w", err)
	}

	key := backupKey(time.Now(), s.cfg.Prefix, s.cfg.DatabaseName)
	if err := s.objectStorage.UploadObject(ctx, key, backup); err != nil {
		return fmt.Errorf("uploading backup: %w", err)
	}

	return nil
}

func (s *Service) PruneDatabaseBackups(ctx context.Context) error {
	keys, err := s.objectStorage.ListObjectKeysPrefix(ctx, s.cfg.Prefix)
	if err != nil {
		return fmt.Errorf("listing backup keys: %w", err)
	}

	threshold := time.Now().Add(-s.cfg.Retention)
	keysToPrune := pickBackupKeysBefore(keys, threshold)

	if err := s.objectStorage.DeleteObjects(ctx, keysToPrune); err != nil {
		return fmt.Errorf("deleting old backups: %w", err)
	}

	return nil
}

func backupTimeString(t time.Time) string {
	return t.UTC().Format("20060102150405")
}

func backupKey(t time.Time, prefix, database string) string {
	filename := fmt.Sprintf("%s-backup-%s.sql", backupTimeString(t), database)
	return path.Join(prefix, filename)
}

func pickBackupKeysBefore(backups []string, t time.Time) []string {
	threshold := backupTimeString(t)

	var picked []string
	for _, name := range backups {
		time := regexPatternDateTime.FindString(name)
		if time != "" && time < threshold {
			picked = append(picked, name)
		}
	}

	slices.Sort(picked)
	return picked
}
