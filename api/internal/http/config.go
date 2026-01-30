package http

import "time"

type Config struct {
	Port                  int           `mapstructure:"port" validate:"required"`
	CookieTokenExpiration time.Duration `mapstructure:"cookie-token-expiration"`
	ShowMetrics           bool          `mapstructure:"show-metrics"`
}
