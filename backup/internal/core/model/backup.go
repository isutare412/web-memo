package model

type DatabaseBackupRequest struct {
	DatabaseName string
	User         string
	Password     string
}
