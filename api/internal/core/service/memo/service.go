package memo

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/port"
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

func (s *Service) GetMemo(ctx context.Context, memoID uuid.UUID) (*ent.Memo, error) {
	memo, err := s.memoRepository.FindByID(ctx, memoID)
	if err != nil {
		return nil, fmt.Errorf("finding memo: %w", err)
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

func (s *Service) CreateMemo(ctx context.Context, memo *ent.Memo, tags []*ent.Tag, userID uuid.UUID) (*ent.Memo, error) {
	tags = sortUniqueTags(tags)

	var memoCreated *ent.Memo

	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		tagsCreated, err := s.ensureTags(ctx, tags)
		if err != nil {
			return fmt.Errorf("ensuring tags: %w", err)
		}

		tagIDs := lo.Map(tagsCreated, func(tag *ent.Tag, _ int) int { return tag.ID })
		memo, err := s.memoRepository.Create(ctx, memo, userID, tagIDs)
		if err != nil {
			return fmt.Errorf("creating memo: %w", err)
		}

		memoCreated = memo
		memoCreated.Edges.Tags = tagsCreated
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("during transaction: %w", err)
	}

	return memoCreated, nil
}

func (s *Service) DeleteMemo(ctx context.Context, memoID uuid.UUID) error {
	if err := s.memoRepository.Delete(ctx, memoID); err != nil {
		return fmt.Errorf("deleting memo: %w", err)
	}
	return nil
}

func (s *Service) ReplaceTags(ctx context.Context, memoID uuid.UUID, tags []*ent.Tag) ([]*ent.Tag, error) {
	tags = sortUniqueTags(tags)

	var tagsReplaced []*ent.Tag

	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		tags, err := s.ensureTags(ctx, tags)
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

func (s *Service) ensureTags(ctx context.Context, tags []*ent.Tag) ([]*ent.Tag, error) {
	tagsCreated := make([]*ent.Tag, 0, len(tags))
	for _, tag := range tags {
		tag, err := s.tagRepository.CreateIfNotExist(ctx, tag.Name)
		if err != nil {
			return nil, fmt.Errorf("creating tag if not exist: %w", err)
		}

		tagsCreated = append(tagsCreated, tag)
	}

	return tagsCreated, nil
}

func sortUniqueTags(tags []*ent.Tag) []*ent.Tag {
	slices.SortStableFunc(tags, func(a, b *ent.Tag) int { return strings.Compare(a.Name, b.Name) })
	return lo.UniqBy(tags, func(tag *ent.Tag) string { return tag.Name })
}
