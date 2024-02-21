package postgres

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/ent/memo"
	"github.com/isutare412/web-memo/api/internal/core/ent/subscription"
	"github.com/isutare412/web-memo/api/internal/core/ent/user"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
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
	client := transactionClient(ctx, r.client)

	users, err := client.User.
		Query().
		Where(user.HasSubscribingMemosWith(memo.ID(memoID))).
		Order(user.BySubscriptions(sql.OrderByField(subscription.FieldCreateTime, sql.OrderDesc()))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Upsert(ctx context.Context, usr *ent.User) (*ent.User, error) {
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
