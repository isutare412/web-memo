package http

import "time"

type Config struct {
	Port                  int
	CookieTokenExpiration time.Duration
	EmbeddingEnabled      bool
}
