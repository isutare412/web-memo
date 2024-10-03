package model

type DatabaseBackupRequest struct {
	Host         string
	Port         int
	DatabaseName string
	User         string
	Password     string
}
