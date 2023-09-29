package postgres

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	QueryLog bool
}
