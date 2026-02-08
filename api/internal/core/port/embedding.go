package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type EmbeddingEnqueuer interface {
	Enqueue(model.EmbeddingJob)
	EnqueueDelete(memoID uuid.UUID)
	Results() <-chan uuid.UUID
}

type EmbeddingSearcher interface {
	Search(ctx context.Context, query string, ownerIDFilter *uuid.UUID, memoIDFilter []uuid.UUID, scoreThreshold float32, limit int) ([]model.SearchResult, error)
}
