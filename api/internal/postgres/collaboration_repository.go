package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/ent/collaboration"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

type CollaborationRepository struct {
	client *ent.Client
}

func NewCollaborationRepository(client *Client) *CollaborationRepository {
	return &CollaborationRepository{
		client: client.entClient,
	}
}

func (r *CollaborationRepository) Find(ctx context.Context, memoID, userID uuid.UUID) (*ent.Collaboration, error) {
	client := transactionClient(ctx, r.client)

	collabo, err := client.Collaboration.
		Query().
		Where(
			collaboration.And(
				collaboration.MemoID(memoID),
				collaboration.UserID(userID),
			),
		).
		Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			ClientMsg: "collaboration not found",
		}
	case err != nil:
		return nil, err
	}

	return collabo, nil
}

func (r *CollaborationRepository) Create(ctx context.Context, memoID, userID uuid.UUID) (*ent.Collaboration, error) {
	client := transactionClient(ctx, r.client)

	clb, err := client.Collaboration.
		Create().
		SetMemoID(memoID).
		SetCollaboratorID(userID).
		SetApproved(false).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return clb, nil
}

func (r *CollaborationRepository) UpdateApprovedStatus(
	ctx context.Context, memoID, userID uuid.UUID, approve bool,
) (*ent.Collaboration, error) {
	client := transactionClient(ctx, r.client)

	count, err := client.Collaboration.
		Update().
		Where(
			collaboration.And(
				collaboration.MemoID(memoID),
				collaboration.UserID(userID),
			),
		).
		SetApproved(approve).
		SetUpdateTime(time.Now()).
		Save(ctx)
	switch {
	case err != nil:
		return nil, err
	case count != 1:
		return nil, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			ClientMsg: "collaboration not found",
		}
	}

	clb, err := client.Collaboration.
		Query().
		Where(
			collaboration.And(
				collaboration.MemoID(memoID),
				collaboration.UserID(userID),
			),
		).
		Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			ClientMsg: "collaboration not found",
		}
	case err != nil:
		return nil, err
	}

	return clb, nil
}

func (r *CollaborationRepository) Delete(ctx context.Context, memoID, userID uuid.UUID) error {
	client := transactionClient(ctx, r.client)

	count, err := client.Collaboration.
		Delete().
		Where(
			collaboration.And(
				collaboration.MemoID(memoID),
				collaboration.UserID(userID),
			),
		).
		Exec(ctx)
	switch {
	case err != nil:
		return err
	case count != 1:
		return pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			ClientMsg: "collaboration not found",
		}
	}

	return nil
}

func (r *CollaborationRepository) DeleteAllByMemoID(
	ctx context.Context,
	memoID uuid.UUID,
) (count int, err error) {
	client := transactionClient(ctx, r.client)

	count, err = client.Collaboration.
		Delete().
		Where(collaboration.MemoID(memoID)).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}
