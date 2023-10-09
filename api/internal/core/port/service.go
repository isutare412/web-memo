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
	ListMemos(ctx context.Context, userID uuid.UUID) ([]*ent.Memo, error)
	CreateMemo(ctx context.Context, memo *ent.Memo, tagNames []string, userID uuid.UUID) (*ent.Memo, error)
	DeleteMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error
	ReplaceTags(ctx context.Context, memoID uuid.UUID, tagNames []string, requester *model.AppIDToken) ([]*ent.Tag, error)
}
