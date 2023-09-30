package postgres

import (
	"context"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/ent/tag"
)

type TagRepository struct {
	client *ent.Client
}

func NewTagRepository(client *Client) *TagRepository {
	return &TagRepository{
		client: client.entClient,
	}
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
