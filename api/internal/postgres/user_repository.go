package postgres

import (
	"context"
	"fmt"
	"slices"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/ent/collaboration"
	"github.com/isutare412/web-memo/api/internal/core/ent/subscription"
	"github.com/isutare412/web-memo/api/internal/core/ent/user"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/tracing"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *Client) *UserRepository {
	return &UserRepository{
		client: client.entClient,
	}
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.UserRepository.FindByID")
	defer span.End()

	client := transactionClient(ctx, r.client)

	userFound, err := client.User.Get(ctx, id)
	switch {
	case ent.IsNotFound(err):
		return nil, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			Origin:    err,
			ClientMsg: fmt.Sprintf("user with id(%s) does not exist", id.String()),
		}
	case err != nil:
		return nil, err
	}

	return userFound, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.UserRepository.FindByEmail")
	defer span.End()

	client := transactionClient(ctx, r.client)

	userFound, err := client.User.
		Query().
		Where(user.Email(email)).
		Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			Origin:    err,
			ClientMsg: fmt.Sprintf("user with email(%s) does not exist", email),
		}
	case err != nil:
		return nil, err
	}

	return userFound, nil
}

func (r *UserRepository) FindAllBySubscribingMemoID(ctx context.Context, memoID uuid.UUID) ([]*ent.User, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.UserRepository.FindAllBySubscribingMemoID")
	defer span.End()

	client := transactionClient(ctx, r.client)

	users, err := client.User.
		Query().
		Where(user.HasSubscriptionsWith(subscription.MemoID(memoID))).
		WithSubscriptions().
		All(ctx)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(users, func(a, b *ent.User) int {
		as, ok := lo.Find(a.Edges.Subscriptions, func(s *ent.Subscription) bool { return s.MemoID == memoID })
		if !ok {
			panic(fmt.Sprintf("subscription not found from user(%v)", a.ID))
		}

		bs, ok := lo.Find(b.Edges.Subscriptions, func(s *ent.Subscription) bool { return s.MemoID == memoID })
		if !ok {
			panic(fmt.Sprintf("subscription not found from user(%v)", b.ID))
		}

		return lo.Ternary(as.CreateTime.After(bs.CreateTime), 1, -1)
	})

	return users, nil
}

func (r *UserRepository) FindAllByCollaboratingMemoIDWithEdges(
	ctx context.Context,
	memoID uuid.UUID,
) ([]*ent.User, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.UserRepository.FindAllByCollaboratingMemoIDWithEdges")
	defer span.End()

	client := transactionClient(ctx, r.client)

	users, err := client.User.
		Query().
		Where(user.HasCollaborationsWith(collaboration.MemoID(memoID))).
		WithCollaborations().
		All(ctx)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(users, func(a, b *ent.User) int {
		as, ok := lo.Find(a.Edges.Subscriptions, func(s *ent.Subscription) bool { return s.MemoID == memoID })
		if !ok {
			panic(fmt.Sprintf("collaboration not found from user(%v)", a.ID))
		}

		bs, ok := lo.Find(b.Edges.Subscriptions, func(s *ent.Subscription) bool { return s.MemoID == memoID })
		if !ok {
			panic(fmt.Sprintf("collaboration not found from user(%v)", b.ID))
		}

		return lo.Ternary(as.CreateTime.After(bs.CreateTime), 1, -1)
	})

	return users, nil
}

func (r *UserRepository) Upsert(ctx context.Context, usr *ent.User) (*ent.User, error) {
	ctx, span := tracing.StartSpan(ctx, "postgres.UserRepository.Upsert")
	defer span.End()

	client := transactionClient(ctx, r.client)

	idUpserted, err := client.User.
		Create().
		SetEmail(usr.Email).
		SetUserName(usr.UserName).
		SetNillableGivenName(lo.EmptyableToPtr(usr.GivenName)).
		SetNillableFamilyName(lo.EmptyableToPtr(usr.FamilyName)).
		SetNillablePhotoURL(lo.EmptyableToPtr(usr.PhotoURL)).
		SetType(usr.Type).
		OnConflict(
			sql.ConflictColumns(user.FieldEmail),
			sql.ResolveWithNewValues(),
			sql.ResolveWith(func(us *sql.UpdateSet) {
				us.SetIgnore(user.FieldID)
				us.SetIgnore(user.FieldCreateTime)

				if usr.GivenName == "" {
					us.SetNull(user.FieldGivenName)
				}
				if usr.FamilyName == "" {
					us.SetNull(user.FieldFamilyName)
				}
				if usr.PhotoURL == "" {
					us.SetNull(user.FieldPhotoURL)
				}
			}),
		).
		ID(ctx)
	if err != nil {
		return nil, err
	}

	userUpserted, err := client.User.Get(ctx, idUpserted)
	if err != nil {
		return nil, err
	}

	return userUpserted, nil
}
