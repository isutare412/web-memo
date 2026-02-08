package port

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/isutare412/imageer/pkg/images"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/model"
)

type AuthService interface {
	VerifyAppIDToken(string) (*model.AppIDToken, error)
	RefreshAppIDToken(ctx context.Context, tokenString string) (newToken *model.AppIDToken, newTokenString string, err error)
	StartGoogleSignIn(context.Context, *http.Request) (redirectURL string, err error)
	FinishGoogleSignIn(context.Context, *http.Request) (redirectURL, appToken string, err error)
}

type MemoService interface {
	GetMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) (*ent.Memo, error)
	SearchMemos(ctx context.Context, userID uuid.UUID, query string) ([]*model.MemoSearchResult, error)
	ListMemos(
		ctx context.Context, userID uuid.UUID, tags []string,
		sortParams model.MemoSortParams, pageParams model.PaginationParams) (memos []*ent.Memo, totalCount int, err error)
	CreateMemo(ctx context.Context, memo *ent.Memo, tagNames []string, userID uuid.UUID) (*ent.Memo, error)
	UpdateMemo(
		ctx context.Context, memo *ent.Memo, tagNames []string,
		requester *model.AppIDToken, isPinUpdateTime bool) (*ent.Memo, error)
	UpdateMemoPublishedState(
		ctx context.Context, memoID uuid.UUID, publish bool, requester *model.AppIDToken) (*ent.Memo, error)
	DeleteMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error

	ListTags(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) ([]*ent.Tag, error)
	SearchTags(ctx context.Context, keyword string, requester *model.AppIDToken) ([]*ent.Tag, error)
	ReplaceTags(ctx context.Context, memoID uuid.UUID, tagNames []string, requester *model.AppIDToken) ([]*ent.Tag, error)
	DeleteOrphanTags(context.Context) (count int, err error)
	EnqueueMissingEmbeddings(context.Context) (enqueued int, err error)

	ListSubscribers(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) (*model.ListSubscribersResponse, error)
	SubscribeMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error
	UnsubscribeMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error

	ListCollaborators(
		ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) (*model.ListCollaboratorsResponse, error)
	RegisterCollaborator(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error
	AuthorizeCollaborator(
		ctx context.Context, memoID, collaboratorID uuid.UUID, approve bool, requester *model.AppIDToken) error
	DeleteCollaborator(ctx context.Context, memoID, collaboratorID uuid.UUID, requester *model.AppIDToken) error
}

type ImageService interface {
	CreateUploadURL(ctx context.Context, fileName string, format images.Format) (*model.UploadURL, error)
	GetImage(ctx context.Context, imageID string, waitUntilProcessed bool) (*model.Image, error)
}
