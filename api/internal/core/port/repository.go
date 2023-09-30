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
