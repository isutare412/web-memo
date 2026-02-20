package port

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
	"github.com/isutare412/web-memo/api/internal/core/model"
)

type TransactionManager interface {
	BeginTx(context.Context) (ctxWithTx context.Context, commit, rollback func() error)
	WithTx(context.Context, func(ctxWithTx context.Context) error) error
}

type UserRepository interface {
	FindByID(ctx context.Context, userID uuid.UUID) (*ent.User, error)
	FindByEmail(ctx context.Context, email string) (*ent.User, error)
	FindAllBySubscribingMemoID(ctx context.Context, memoID uuid.UUID) ([]*ent.User, error)
	FindAllByCollaboratingMemoIDWithEdges(ctx context.Context, memoID uuid.UUID) ([]*ent.User, error)
	Upsert(context.Context, *ent.User) (*ent.User, error)
}

type MemoRepository interface {
	FindByID(ctx context.Context, memoID uuid.UUID) (*ent.Memo, error)
	FindByIDWithEdges(ctx context.Context, memoID uuid.UUID) (*ent.Memo, error)
	FindSubscribedMemoIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	FindAllNotEmbedded(ctx context.Context, pageParams model.PaginationParams) ([]*ent.Memo, error)
	FindAll(ctx context.Context, sortParams model.MemoSortParams, pageParams model.PaginationParams) ([]*ent.Memo, error)
	FindAllByUserIDWithEdges(
		ctx context.Context, userID uuid.UUID,
		sortParams model.MemoSortParams, pageParams model.PaginationParams) ([]*ent.Memo, error)
	FindAllByIDsWithEdges(ctx context.Context, ids []uuid.UUID) ([]*ent.Memo, error)
	FindAllByUserIDAndTagNamesWithEdges(
		ctx context.Context, userID uuid.UUID, tags []string,
		sortParams model.MemoSortParams, pageParams model.PaginationParams) ([]*ent.Memo, error)
	CountByUserIDAndTagNames(ctx context.Context, userID uuid.UUID, tags []string) (int, error)
	Create(ctx context.Context, memo *ent.Memo, userID uuid.UUID, tagIDs []int) (*ent.Memo, error)
	Update(context.Context, *ent.Memo) (*ent.Memo, error)
	UpdatePublishState(ctx context.Context, memoID uuid.UUID, state enum.PublishState) (*ent.Memo, error)
	UpdateIsEmbedded(ctx context.Context, memoID uuid.UUID, isEmbedded bool) error
	Delete(ctx context.Context, memoID uuid.UUID) error

	ReplaceTags(ctx context.Context, memoID uuid.UUID, tagIDs []int, updateTime bool) error

	FindSubscription(ctx context.Context, memoID, userID uuid.UUID) (*ent.Subscription, error)
	FindSubscriptionsByMemoID(ctx context.Context, memoID uuid.UUID) ([]*ent.Subscription, error)
	RegisterSubscriber(ctx context.Context, memoID, userID uuid.UUID, approved bool) error
	UpdateSubscriptionApproval(ctx context.Context, memoID, userID uuid.UUID, approved bool) error
	ApproveAllSubscriptions(ctx context.Context, memoID uuid.UUID) error
	UnregisterSubscriber(ctx context.Context, memoID, userID uuid.UUID) error
	ClearSubscribers(ctx context.Context, memoID uuid.UUID) error
}

type TagRepository interface {
	FindAllByMemoID(ctx context.Context, memoID uuid.UUID) ([]*ent.Tag, error)
	FindAllByUserIDAndNameContains(ctx context.Context, userID uuid.UUID, name string) ([]*ent.Tag, error)
	CreateIfNotExist(ctx context.Context, tagName string) (*ent.Tag, error)
	DeleteAllWithoutMemo(ctx context.Context, excludes []string) (count int, err error)
}

type CollaborationRepository interface {
	Find(ctx context.Context, memoID, userID uuid.UUID) (*ent.Collaboration, error)
	Create(ctx context.Context, memoID, userID uuid.UUID) (*ent.Collaboration, error)
	UpdateApprovedStatus(ctx context.Context, memoID, userID uuid.UUID, approve bool) (*ent.Collaboration, error)
	Delete(ctx context.Context, memoID, userID uuid.UUID) error
	DeleteAllByMemoID(
		ctx context.Context,
		memoID uuid.UUID,
	) (count int, err error)
}

type KVRepository interface {
	Get(ctx context.Context, key string) (string, error)
	GetThenDelete(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, val string, exp time.Duration) error
	Delete(ctx context.Context, keys ...string) (delCount int64, err error)
}
