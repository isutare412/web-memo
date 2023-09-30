package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/ent"
)

type MemoService interface {
	GetMemo(ctx context.Context, memoID uuid.UUID) (*ent.Memo, error)
	ListMemos(ctx context.Context, userID uuid.UUID) ([]*ent.Memo, error)
	CreateMemo(ctx context.Context, memo *ent.Memo, tags []*ent.Tag, userID uuid.UUID) (*ent.Memo, error)
	DeleteMemo(ctx context.Context, memoID uuid.UUID) error
	ReplaceTags(ctx context.Context, memoID uuid.UUID, tags []*ent.Tag) ([]*ent.Tag, error)
}
