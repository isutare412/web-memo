package port

import (
	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type EmbeddingEnqueuer interface {
	Enqueue(model.EmbeddingJob)
	EnqueueDelete(memoID uuid.UUID)
	Results() <-chan uuid.UUID
}
