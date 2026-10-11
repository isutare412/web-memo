package postgres

type Config struct {
	Host     string `koanf:"host" validate:"required"`
	Port     int    `koanf:"port" validate:"required"`
	User     string `koanf:"user" validate:"required"`
	Password string `koanf:"password" validate:"required"`
	Database string `koanf:"database" validate:"required"`
	QueryLog bool   `koanf:"query-log"`
	// SSLMode is the libpq sslmode. Empty means disable (local compose has no TLS).
	SSLMode string `koanf:"ssl-mode" validate:"omitempty,oneof=disable allow prefer require"`
}
