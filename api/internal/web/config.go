package web

import "time"

type Config struct {
	Port                  int
	ShowOpenAPIDocs       bool
	CookieTokenExpiration time.Duration
	EmbeddingEnabled      bool
}
