package port

import "context"

type BackupService interface {
	BackupDatabase(context.Context) error
	PruneDatabaseBackups(context.Context) error
}
