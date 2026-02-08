package embedding

import (
	"context"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type NoopEnqueuer struct{}

func (NoopEnqueuer) Enqueue(model.EmbeddingJob) {}
func (NoopEnqueuer) EnqueueDelete(uuid.UUID)    {}
func (NoopEnqueuer) Results() <-chan uuid.UUID  { return nil }

type NoopSearcher struct{}

func (NoopSearcher) Search(context.Context, string, *uuid.UUID, []uuid.UUID, float32, int) ([]model.SearchResult, error) {
	return nil, nil
}
