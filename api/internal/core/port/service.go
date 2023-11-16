package port

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/model"
)

type AuthService interface {
	VerifyAppIDTokenString(string) (*model.AppIDToken, error)
	StartGoogleSignIn(context.Context, *http.Request) (redirectURL string, err error)
	FinishGoogleSignIn(context.Context, *http.Request) (redirectURL, appToken string, err error)
}

type MemoService interface {
	GetMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) (*ent.Memo, error)
	ListMemos(ctx context.Context, userID uuid.UUID, option *model.QueryOption) (memos []*ent.Memo, totalCount int, err error)
	CreateMemo(ctx context.Context, memo *ent.Memo, tagNames []string, userID uuid.UUID) (*ent.Memo, error)
	UpdateMemo(ctx context.Context, memo *ent.Memo, tagNames []string, requester *model.AppIDToken) (*ent.Memo, error)
	DeleteMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error

	ListTags(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) ([]*ent.Tag, error)
	SearchTags(ctx context.Context, keyword string, requester *model.AppIDToken) ([]*ent.Tag, error)
	ReplaceTags(ctx context.Context, memoID uuid.UUID, tagNames []string, requester *model.AppIDToken) ([]*ent.Tag, error)
	DeleteOrphanTags(context.Context) (count int, err error)
}
