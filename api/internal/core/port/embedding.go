package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type EmbeddingEnqueuer interface {
	Enqueue(model.EmbeddingJob)
	EnqueueDelete(memoID uuid.UUID)
}

type EmbeddingRepository interface {
	ExistsByMemoID(ctx context.Context, memoID uuid.UUID) (bool, error)
}
