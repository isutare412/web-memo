package model

import "github.com/google/uuid"

type EmbeddingJob struct {
	MemoID  uuid.UUID
	OwnerID uuid.UUID
	Title   string
	Content string
}

type SearchResult struct {
	MemoID uuid.UUID
	Score  float32
}
