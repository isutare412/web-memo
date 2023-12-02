package postgres

import (
	"context"
	"encoding/base64"
	"fmt"
	"slices"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/ent/memo"
	"github.com/isutare412/web-memo/api/internal/core/ent/predicate"
	"github.com/isutare412/web-memo/api/internal/core/ent/tag"
	"github.com/isutare412/web-memo/api/internal/core/ent/user"
	"github.com/isutare412/web-memo/api/internal/core/model"
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

	contentDecoded, err := base64Decode(memo.Content)
	if err != nil {
		return nil, err
	}
	memo.Content = contentDecoded

	return memo, nil
}

func (r *MemoRepository) FindByIDWithTags(ctx context.Context, memoID uuid.UUID) (*ent.Memo, error) {
	client := transactionClient(ctx, r.client)

	memo, err := client.Memo.
		Query().
		Where(memo.ID(memoID)).
		WithTags().
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

	slices.SortFunc(memo.Edges.Tags, func(a, b *ent.Tag) int { return strings.Compare(a.Name, b.Name) })

	contentDecoded, err := base64Decode(memo.Content)
	if err != nil {
		return nil, err
	}
	memo.Content = contentDecoded

	return memo, nil
}

func (r *MemoRepository) FindAllByUserIDWithTags(
	ctx context.Context,
	userID uuid.UUID,
	opt *model.QueryOption,
) ([]*ent.Memo, error) {
	client := transactionClient(ctx, r.client)

	createTimeOrder := memo.ByCreateTime(sql.OrderDesc())
	if opt.Direction == model.SortDirectionAsc {
		createTimeOrder = memo.ByCreateTime(sql.OrderAsc())
	}

	query := client.Memo.
		Query().
		Where(memo.HasOwnerWith(user.ID(userID))).
		WithTags().
		Order(createTimeOrder)

	if opt.PageOffset > 0 && opt.PageSize > 0 {
		offset := (opt.PageOffset - 1) * opt.PageSize
		query = query.Offset(offset).Limit(opt.PageSize)
	}

	memos, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	for _, m := range memos {
		slices.SortFunc(m.Edges.Tags, func(a, b *ent.Tag) int { return strings.Compare(a.Name, b.Name) })

		contentDecoded, err := base64Decode(m.Content)
		if err != nil {
			return nil, err
		}
		m.Content = contentDecoded
	}

	return memos, nil
}

func (r *MemoRepository) FindAllByUserIDAndTagNamesWithTags(
	ctx context.Context,
	userID uuid.UUID,
	tags []string,
	opt *model.QueryOption,
) ([]*ent.Memo, error) {
	client := transactionClient(ctx, r.client)

	createTimeOrder := memo.ByCreateTime(sql.OrderDesc())
	if opt.Direction == model.SortDirectionAsc {
		createTimeOrder = memo.ByCreateTime(sql.OrderAsc())
	}

	query := client.Memo.
		Query().
		Where(memo.HasOwnerWith(user.ID(userID))).
		WithTags().
		Order(createTimeOrder)

	if opt.PageOffset > 0 && opt.PageSize > 0 {
		offset := (opt.PageOffset - 1) * opt.PageSize
		query = query.Offset(offset).Limit(opt.PageSize)
	}

	if len(tags) > 0 {
		tagsMatch := lo.Map(
			tags,
			func(tagName string, _ int) predicate.Memo { return memo.HasTagsWith(tag.Name(tagName)) })
		query = query.Where(memo.And(tagsMatch...))
	}

	memos, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	for _, m := range memos {
		slices.SortFunc(m.Edges.Tags, func(a, b *ent.Tag) int { return strings.Compare(a.Name, b.Name) })

		contentDecoded, err := base64Decode(m.Content)
		if err != nil {
			return nil, err
		}
		m.Content = contentDecoded
	}

	return memos, nil
}

func (r *MemoRepository) CountByUserIDAndTagNames(ctx context.Context, userID uuid.UUID, tags []string) (int, error) {
	client := transactionClient(ctx, r.client)

	query := client.Memo.
		Query().
		Where(memo.OwnerID(userID))

	if len(tags) > 0 {
		tagsMatch := lo.Map(
			tags,
			func(tagName string, _ int) predicate.Memo { return memo.HasTagsWith(tag.Name(tagName)) })
		query = query.Where(memo.And(tagsMatch...))
	}

	count, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
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
		SetContent(base64Encode(memo.Content)).
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
		SetContent(base64Encode(memo.Content)).
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

func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func base64Decode(s string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("base64 decoding string: %w", err)
	}
	return string(decodedBytes), nil
}
