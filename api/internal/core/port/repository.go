package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/ent"
)

type UserRepository interface {
	FindByID(context.Context, uuid.UUID) (*ent.User, error)
	Upsert(context.Context, *ent.User) (*ent.User, error)
}

type MemoRepository interface {
	FindAllByUserIDWithTags(context.Context, uuid.UUID) ([]*ent.Memo, error)
	FindAllByUserIDAndTagIDWithTags(ctx context.Context, userID uuid.UUID, tagID int) ([]*ent.Memo, error)
	Create(ctx context.Context, memo *ent.Memo, userID uuid.UUID, tagIDs []int) (*ent.Memo, error)
	Update(context.Context, *ent.Memo) (*ent.Memo, error)
	ReplaceTags(ctx context.Context, memoID uuid.UUID, tagIDs []int) error
}

type TagRepository interface {
	CreateIfNotExist(ctx context.Context, tagName string) (*ent.Tag, error)
}
