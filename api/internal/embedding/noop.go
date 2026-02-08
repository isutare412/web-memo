package embedding

import (
	"context"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/model"
)

type NoopEnqueuer struct{}

func (NoopEnqueuer) Enqueue(model.EmbeddingJob) {}
func (NoopEnqueuer) EnqueueDelete(uuid.UUID)    {}

type NoopRepository struct{}

func (NoopRepository) ExistsByMemoID(context.Context, uuid.UUID) (bool, error) { return true, nil }
