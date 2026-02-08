package cron

import "time"

type Config struct {
	TagCleanupInterval    time.Duration
	EmbeddingSyncInterval time.Duration
	EmbeddingSyncEnabled  bool
}
