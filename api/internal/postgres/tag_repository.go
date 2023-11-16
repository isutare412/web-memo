package postgres

import (
	"context"
	"slices"
	"strings"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/ent/memo"
	"github.com/isutare412/web-memo/api/internal/core/ent/tag"
	"github.com/isutare412/web-memo/api/internal/core/ent/user"
)

type TagRepository struct {
	client *ent.Client
}

func NewTagRepository(client *Client) *TagRepository {
	return &TagRepository{
		client: client.entClient,
	}
}

func (r *TagRepository) FindAllByMemoID(ctx context.Context, memoID uuid.UUID) ([]*ent.Tag, error) {
	client := transactionClient(ctx, r.client)

	tags, err := client.Tag.
		Query().
		Where(tag.HasMemosWith(memo.ID(memoID))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(tags, func(a, b *ent.Tag) int { return strings.Compare(a.Name, b.Name) })

	return tags, nil
}

func (r *TagRepository) FindAllByUserIDAndNameContains(
	ctx context.Context,
	userID uuid.UUID,
	name string,
) ([]*ent.Tag, error) {
	client := transactionClient(ctx, r.client)

	tags, err := client.Tag.
		Query().
		Where(
			tag.HasMemosWith(memo.HasOwnerWith(user.ID(userID))),
			tag.NameContainsFold(name),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(tags, func(a, b *ent.Tag) int { return strings.Compare(a.Name, b.Name) })

	return tags, nil
}

func (r *TagRepository) CreateIfNotExist(ctx context.Context, tagName string) (*ent.Tag, error) {
	client := transactionClient(ctx, r.client)

	tagFound, err := client.Tag.
		Query().
		Where(tag.Name(tagName)).
		First(ctx)
	if err == nil {
		return tagFound, nil
	}

	tagCreated, err := client.Tag.
		Create().
		SetName(tagName).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return tagCreated, nil
}

func (r *TagRepository) DeleteAllWithoutMemo(ctx context.Context) (count int, err error) {
	client := transactionClient(ctx, r.client)

	count, err = client.Tag.
		Delete().
		Where(tag.Not(tag.HasMemos())).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}
