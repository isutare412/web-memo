package memo

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

type Service struct {
	transactionManager port.TransactionManager
	memoRepository     port.MemoRepository
	tagRepository      port.TagRepository
}

func NewService(
	transactionManager port.TransactionManager,
	memoRepository port.MemoRepository,
	tagRepository port.TagRepository,
) *Service {
	return &Service{
		transactionManager: transactionManager,
		memoRepository:     memoRepository,
		tagRepository:      tagRepository,
	}
}

func (s *Service) GetMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) (*ent.Memo, error) {
	memo, err := s.memoRepository.FindByIDWithTags(ctx, memoID)
	if err != nil {
		return nil, fmt.Errorf("finding memo: %w", err)
	}

	if !requester.CanAccessMemo(memo) {
		return nil, pkgerr.Known{Code: pkgerr.CodePermissionDenied}
	}

	return memo, nil
}

func (s *Service) ListMemos(ctx context.Context, userID uuid.UUID) ([]*ent.Memo, error) {
	memos, err := s.memoRepository.FindAllByUserIDWithTags(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("finding all memos of user(%s): %w", userID.String(), err)
	}

	return memos, nil
}

func (s *Service) CreateMemo(ctx context.Context, memo *ent.Memo, tagNames []string, userID uuid.UUID) (*ent.Memo, error) {
	tagNames = sortDedupTags(tagNames)

	var memoCreated *ent.Memo

	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		tags, err := s.ensureTags(ctx, tagNames)
		if err != nil {
			return fmt.Errorf("ensuring tags: %w", err)
		}

		tagIDs := lo.Map(tags, func(tag *ent.Tag, _ int) int { return tag.ID })
		memo, err := s.memoRepository.Create(ctx, memo, userID, tagIDs)
		if err != nil {
			return fmt.Errorf("creating memo: %w", err)
		}

		memoCreated = memo
		memoCreated.Edges.Tags = tags
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("during transaction: %w", err)
	}

	return memoCreated, nil
}

func (s *Service) UpdateMemo(
	ctx context.Context,
	memo *ent.Memo,
	tagNames []string,
	requester *model.AppIDToken,
) (*ent.Memo, error) {
	tagNames = sortDedupTags(tagNames)

	var memoUpdated *ent.Memo

	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		tags, err := s.ensureTags(ctx, tagNames)
		if err != nil {
			return fmt.Errorf("ensuring tags: %w", err)
		}

		memoFound, err := s.memoRepository.FindByID(ctx, memo.ID)
		if err != nil {
			return fmt.Errorf("finding memo: %w", err)
		}

		if !requester.CanAccessMemo(memoFound) {
			return pkgerr.Known{Code: pkgerr.CodePermissionDenied}
		}

		memo, err := s.memoRepository.Update(ctx, memo)
		if err != nil {
			return fmt.Errorf("creating memo: %w", err)
		}

		tagIDs := lo.Map(tags, func(tag *ent.Tag, _ int) int { return tag.ID })
		if err := s.memoRepository.ReplaceTags(ctx, memo.ID, tagIDs); err != nil {
			return fmt.Errorf("replacing tags: %w", err)
		}

		memoUpdated = memo
		memoUpdated.Edges.Tags = tags
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("during transaction: %w", err)
	}

	return memoUpdated, nil
}

func (s *Service) DeleteMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error {
	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		memo, err := s.memoRepository.FindByID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding memo: %w", err)
		}

		if !requester.CanAccessMemo(memo) {
			return pkgerr.Known{Code: pkgerr.CodePermissionDenied}
		}

		if err := s.memoRepository.Delete(ctx, memoID); err != nil {
			return fmt.Errorf("deleting memo: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("during transaction: %w", err)
	}

	return nil
}

func (s *Service) ListTags(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) ([]*ent.Tag, error) {
	var tagsFound []*ent.Tag

	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		memo, err := s.memoRepository.FindByID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding memo: %w", err)
		}

		if !requester.CanAccessMemo(memo) {
			return pkgerr.Known{Code: pkgerr.CodePermissionDenied}
		}

		tags, err := s.tagRepository.FindAllByMemoID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding all tags: %w", err)
		}

		tagsFound = tags
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("during transaction: %w", err)
	}

	return tagsFound, nil
}

func (s *Service) ReplaceTags(
	ctx context.Context,
	memoID uuid.UUID,
	tagNames []string,
	requester *model.AppIDToken,
) ([]*ent.Tag, error) {
	tagNames = sortDedupTags(tagNames)

	var tagsReplaced []*ent.Tag

	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		memo, err := s.memoRepository.FindByID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding memo: %w", err)
		}

		if !requester.CanAccessMemo(memo) {
			return pkgerr.Known{Code: pkgerr.CodePermissionDenied}
		}

		tags, err := s.ensureTags(ctx, tagNames)
		if err != nil {
			return fmt.Errorf("ensuring tags: %w", err)
		}

		tagIDs := lo.Map(tags, func(tag *ent.Tag, _ int) int { return tag.ID })
		if err := s.memoRepository.ReplaceTags(ctx, memoID, tagIDs); err != nil {
			return fmt.Errorf("replacing tags: %w", err)
		}

		tagsReplaced = tags
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("during transaction: %w", err)
	}

	return tagsReplaced, nil
}

func (s *Service) ensureTags(ctx context.Context, tagNames []string) ([]*ent.Tag, error) {
	if _, ok := lo.Find(tagNames, func(tag string) bool { return len(tag) > 20 }); ok {
		return nil, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: fmt.Sprintf("length of tag should be less or equal to 20"),
		}
	}

	tagsCreated := make([]*ent.Tag, 0, len(tagNames))
	for _, tagName := range tagNames {
		tag, err := s.tagRepository.CreateIfNotExist(ctx, tagName)
		if err != nil {
			return nil, fmt.Errorf("creating tag if not exist: %w", err)
		}

		tagsCreated = append(tagsCreated, tag)
	}

	return tagsCreated, nil
}

func sortDedupTags(tags []string) []string {
	tags = lo.Uniq(tags)
	slices.Sort(tags)
	return tags
}
