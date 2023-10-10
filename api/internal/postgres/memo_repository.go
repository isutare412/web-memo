package postgres

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/ent/memo"
	"github.com/isutare412/web-memo/api/internal/core/ent/tag"
	"github.com/isutare412/web-memo/api/internal/core/ent/user"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

type MemoRepository struct {
	client *ent.Client
}

func NewMemoRepository(client *Client) *MemoRepository {
	return &MemoRepository{
		client: client.entClient,
	}
}

func (r *MemoRepository) FindByID(ctx context.Context, memoID uuid.UUID) (*ent.Memo, error) {
	client := transactionClient(ctx, r.client)

	memo, err := client.Memo.
		Query().
		Where(memo.ID(memoID)).
		First(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			Origin:    err,
			ClientMsg: fmt.Sprintf("memo with id(%s) not found", memoID.String()),
		}
	case err != nil:
		return nil, err
	}

	return memo, nil
}

func (r *MemoRepository) FindByIDWithTags(ctx context.Context, memoID uuid.UUID) (*ent.Memo, error) {
	client := transactionClient(ctx, r.client)

	memo, err := client.Memo.
		Query().
		Where(memo.ID(memoID)).
		WithTags(func(tq *ent.TagQuery) {
			tq.Order(tag.ByName(sql.OrderAsc()))
		}).
		First(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			Origin:    err,
			ClientMsg: fmt.Sprintf("memo with id(%s) not found", memoID.String()),
		}
	case err != nil:
		return nil, err
	}

	return memo, nil
}

func (r *MemoRepository) FindAllByUserIDWithTags(ctx context.Context, userID uuid.UUID) ([]*ent.Memo, error) {
	client := transactionClient(ctx, r.client)

	memos, err := client.Memo.
		Query().
		Where(memo.HasOwnerWith(user.ID(userID))).
		WithTags(func(tq *ent.TagQuery) {
			tq.Order(tag.ByName(sql.OrderAsc()))
		}).
		Order(memo.ByCreateTime(sql.OrderDesc())).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return memos, nil
}

func (r *MemoRepository) FindAllByUserIDAndTagIDWithTags(
	ctx context.Context,
	userID uuid.UUID,
	tagID int,
) ([]*ent.Memo, error) {
	client := transactionClient(ctx, r.client)

	memos, err := client.Memo.
		Query().
		Where(
			memo.And(
				memo.HasOwnerWith(user.ID(userID)),
				memo.HasTagsWith(tag.ID(tagID)),
			),
		).
		WithTags(func(tq *ent.TagQuery) {
			tq.Order(tag.ByName(sql.OrderAsc()))
		}).
		Order(memo.ByCreateTime(sql.OrderDesc())).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return memos, nil
}

func (r *MemoRepository) Create(
	ctx context.Context,
	memo *ent.Memo,
	userID uuid.UUID,
	tagIDs []int,
) (*ent.Memo, error) {
	client := transactionClient(ctx, r.client)

	memoCreated, err := client.Memo.
		Create().
		SetTitle(memo.Title).
		SetContent(memo.Content).
		SetOwnerID(userID).
		AddTagIDs(tagIDs...).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return memoCreated, nil
}

func (r *MemoRepository) Update(ctx context.Context, memo *ent.Memo) (*ent.Memo, error) {
	client := transactionClient(ctx, r.client)

	memoUpdated, err := client.Memo.
		UpdateOneID(memo.ID).
		SetTitle(memo.Title).
		SetContent(memo.Content).
		Save(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			Origin:    err,
			ClientMsg: fmt.Sprintf("memo with id(%s) not found", memo.ID.String()),
		}
	case err != nil:
		return nil, err
	}

	return memoUpdated, nil
}

func (r *MemoRepository) Delete(ctx context.Context, memoID uuid.UUID) error {
	client := transactionClient(ctx, r.client)

	err := client.Memo.
		DeleteOneID(memoID).
		Exec(ctx)
	switch {
	case ent.IsNotFound(err):
		return pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			Origin:    err,
			ClientMsg: fmt.Sprintf("memo with id(%s) does not exist", memoID.String()),
		}
	case err != nil:
		return err
	}

	return nil
}

func (r *MemoRepository) ReplaceTags(ctx context.Context, memoID uuid.UUID, tagIDs []int) error {
	client := transactionClient(ctx, r.client)

	err := client.Memo.
		UpdateOneID(memoID).
		ClearTags().
		AddTagIDs(tagIDs...).
		Exec(ctx)
	switch {
	case ent.IsNotFound(err):
		return pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			Origin:    err,
			ClientMsg: fmt.Sprintf("memo with id(%s) not found", memoID.String()),
		}
	case err != nil:
		return err
	}

	return nil
}
