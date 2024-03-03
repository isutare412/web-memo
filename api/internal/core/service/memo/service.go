package memo

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

const (
	tagPublished = "published"
)

var reservedTags = []string{
	tagPublished,
}

type Service struct {
	transactionManager port.TransactionManager
	memoRepository     port.MemoRepository
	tagRepository      port.TagRepository
	userRepository     port.UserRepository
}

func NewService(
	transactionManager port.TransactionManager,
	memoRepository port.MemoRepository,
	tagRepository port.TagRepository,
	userReposiotry port.UserRepository,
) *Service {
	return &Service{
		transactionManager: transactionManager,
		memoRepository:     memoRepository,
		tagRepository:      tagRepository,
		userRepository:     userReposiotry,
	}
}

func (s *Service) GetMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) (*ent.Memo, error) {
	memo, err := s.memoRepository.FindByIDWithTags(ctx, memoID)
	if err != nil {
		return nil, fmt.Errorf("finding memo: %w", err)
	}

	if !requester.CanReadMemo(memo) {
		return nil, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "not allowed to access memo",
		}
	}

	return memo, nil
}

func (s *Service) ListMemos(
	ctx context.Context,
	userID uuid.UUID,
	tags []string,
	sortParams model.MemoSortParams,
	pageParams model.PaginationParams,
) (memos []*ent.Memo, totalCount int, err error) {
	var (
		memosFound []*ent.Memo
		memoCount  int
	)

	err = s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		memos, err := s.memoRepository.FindAllByUserIDAndTagNamesWithTags(ctx, userID, tags, sortParams, pageParams)
		if err != nil {
			return fmt.Errorf("finding memos of user(%s): %w", userID.String(), err)
		}

		count, err := s.memoRepository.CountByUserIDAndTagNames(ctx, userID, tags)
		if err != nil {
			return fmt.Errorf("getting memo count of user(%s): %w", userID.String(), err)
		}

		memosFound = memos
		memoCount = count
		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("during transaction: %w", err)
	}

	return memosFound, memoCount, nil
}

func (s *Service) CreateMemo(ctx context.Context, memo *ent.Memo, tagNames []string, userID uuid.UUID) (*ent.Memo, error) {
	tagNames = removeReservedTags(tagNames)
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
	isPinUpdateTime bool,
) (*ent.Memo, error) {
	tagNames = removeReservedTags(tagNames)
	tagNames = sortDedupTags(tagNames)

	var memoUpdated *ent.Memo

	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		memoFound, err := s.memoRepository.FindByIDWithTags(ctx, memo.ID)
		if err != nil {
			return fmt.Errorf("finding memo: %w", err)
		}

		if !requester.CanWriteMemo(memoFound) {
			return pkgerr.Known{
				Code:      pkgerr.CodePermissionDenied,
				ClientMsg: "not allowed to access memo",
			}
		}

		tags, err := s.ensureTags(ctx, tagNames)
		if err != nil {
			return fmt.Errorf("ensuring tags: %w", err)
		}

		reservedTags := lo.Filter(
			memoFound.Edges.Tags,
			func(t *ent.Tag, _ int) bool { return lo.Contains(reservedTags, t.Name) })
		tags = append(tags, reservedTags...)

		memo.IsPublished = memoFound.IsPublished
		if isPinUpdateTime {
			memo.UpdateTime = memoFound.UpdateTime
		}

		memo, err := s.memoRepository.Update(ctx, memo)
		if err != nil {
			return fmt.Errorf("creating memo: %w", err)
		}

		tagIDs := lo.Map(tags, func(tag *ent.Tag, _ int) int { return tag.ID })
		if err := s.memoRepository.ReplaceTags(ctx, memo.ID, tagIDs, !isPinUpdateTime); err != nil {
			return fmt.Errorf("replacing tags: %w", err)
		}

		tags, err = s.tagRepository.FindAllByMemoID(ctx, memo.ID)
		if err != nil {
			return fmt.Errorf("finding tags of memo: %w", err)
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

func (s *Service) UpdateMemoPublishedState(
	ctx context.Context,
	memoID uuid.UUID,
	publish bool,
	requester *model.AppIDToken,
) (*ent.Memo, error) {
	var memoUpdated *ent.Memo

	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		memoFound, err := s.memoRepository.FindByIDWithTags(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding memo: %w", err)
		}

		if !requester.CanWriteMemo(memoFound) {
			return pkgerr.Known{
				Code:      pkgerr.CodePermissionDenied,
				ClientMsg: "not allowed to access memo",
			}
		}

		if memoFound.IsPublished == publish {
			memoUpdated = memoFound
			return nil
		}

		var tagIDs []int
		if publish {
			tags, err := s.ensureTags(ctx, []string{tagPublished})
			if err != nil {
				return fmt.Errorf("ensuring tags: %w", err)
			}

			tagIDs = lo.Map(append(memoFound.Edges.Tags, tags...), func(tag *ent.Tag, _ int) int { return tag.ID })
		} else {
			tags := lo.Filter(memoFound.Edges.Tags, func(t *ent.Tag, _ int) bool { return t.Name != tagPublished })
			tagIDs = lo.Map(tags, func(tag *ent.Tag, _ int) int { return tag.ID })
		}

		if err := s.memoRepository.ReplaceTags(ctx, memoFound.ID, tagIDs, false); err != nil {
			return fmt.Errorf("replacing tags: %w", err)
		}

		memo, err := s.memoRepository.UpdateIsPublish(ctx, memoFound.ID, publish)
		if err != nil {
			return fmt.Errorf("updating memo published state: %w", err)
		}

		if !publish {
			if err := s.memoRepository.ClearSubscribers(ctx, memoFound.ID); err != nil {
				return fmt.Errorf("clearing subscribers: %w", err)
			}
		}

		tags, err := s.tagRepository.FindAllByMemoID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding tags of memo: %w", err)
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

		if !requester.CanWriteMemo(memo) {
			return pkgerr.Known{
				Code:      pkgerr.CodePermissionDenied,
				ClientMsg: "not allowed to access memo",
			}
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

		if !requester.CanReadMemo(memo) {
			return pkgerr.Known{
				Code:      pkgerr.CodePermissionDenied,
				ClientMsg: "not allowed to access memo",
			}
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

func (s *Service) SearchTags(ctx context.Context, keyword string, requester *model.AppIDToken) ([]*ent.Tag, error) {
	tags, err := s.tagRepository.FindAllByUserIDAndNameContains(ctx, requester.UserID, keyword)
	if err != nil {
		return nil, fmt.Errorf("finding tags: %w", err)
	}

	return tags, nil
}

func (s *Service) ReplaceTags(
	ctx context.Context,
	memoID uuid.UUID,
	tagNames []string,
	requester *model.AppIDToken,
) ([]*ent.Tag, error) {
	tagNames = removeReservedTags(tagNames)
	tagNames = sortDedupTags(tagNames)

	var tagsReplaced []*ent.Tag

	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		memo, err := s.memoRepository.FindByIDWithTags(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding memo: %w", err)
		}

		if !requester.CanWriteMemo(memo) {
			return pkgerr.Known{
				Code:      pkgerr.CodePermissionDenied,
				ClientMsg: "not allowed to access memo",
			}
		}

		tags, err := s.ensureTags(ctx, tagNames)
		if err != nil {
			return fmt.Errorf("ensuring tags: %w", err)
		}

		reservedTags := lo.Filter(
			memo.Edges.Tags,
			func(t *ent.Tag, _ int) bool { return lo.Contains(reservedTags, t.Name) })
		tags = append(tags, reservedTags...)

		tagIDs := lo.Map(tags, func(tag *ent.Tag, _ int) int { return tag.ID })
		if err := s.memoRepository.ReplaceTags(ctx, memoID, tagIDs, true); err != nil {
			return fmt.Errorf("replacing tags: %w", err)
		}

		tags, err = s.tagRepository.FindAllByMemoID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding tags of memo: %w", err)
		}

		tagsReplaced = tags
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("during transaction: %w", err)
	}

	return tagsReplaced, nil
}

func (s *Service) DeleteOrphanTags(ctx context.Context) (count int, err error) {
	count, err = s.tagRepository.DeleteAllWithoutMemo(ctx, reservedTags)
	if err != nil {
		return 0, fmt.Errorf("deleting all tags without memo: %w", err)
	}

	return count, nil
}

func (s *Service) ListSubscribers(
	ctx context.Context,
	memoID uuid.UUID,
	requester *model.AppIDToken,
) (*model.ListSubscribersResponse, error) {
	var (
		ownerID     uuid.UUID
		subscribers []*ent.User
	)
	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		memo, err := s.memoRepository.FindByID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding memo: %w", err)
		}

		if !requester.CanReadMemo(memo) {
			return pkgerr.Known{
				Code:      pkgerr.CodePermissionDenied,
				ClientMsg: "not allowed to access memo",
			}
		}

		users, err := s.userRepository.FindAllBySubscribingMemoID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding subscribers: %w", err)
		}

		ownerID = memo.OwnerID
		subscribers = users
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("during transaction: %w", err)
	}

	return &model.ListSubscribersResponse{
		MemoOwnerID: ownerID,
		Subscribers: subscribers,
	}, nil
}

func (s *Service) SubscribeMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error {
	if requester == nil {
		return pkgerr.Known{
			Code:      pkgerr.CodeUnauthenticated,
			ClientMsg: "must be signed-in for subscription",
		}
	}

	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		memo, err := s.memoRepository.FindByID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding memo: %w", err)
		}
		if memo.OwnerID == requester.UserID {
			return pkgerr.Known{
				Code:      pkgerr.CodeBadRequest,
				ClientMsg: "cannot subscribe memo of your own",
			}
		}

		if !memo.IsPublished {
			return pkgerr.Known{
				Code:      pkgerr.CodeBadRequest,
				ClientMsg: "memo is not published",
			}
		}

		subscribers, err := s.userRepository.FindAllBySubscribingMemoID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding existing subscribers: %w", err)
		}

		if _, ok := lo.Find(subscribers, func(u *ent.User) bool { return u.ID == requester.UserID }); ok {
			return pkgerr.Known{
				Code:      pkgerr.CodeBadRequest,
				ClientMsg: "already subscribed memo",
			}
		}

		if err := s.memoRepository.RegisterSubscriber(ctx, memoID, requester.UserID); err != nil {
			return fmt.Errorf("registering subscriber: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("during transaction: %w", err)
	}

	return nil
}

func (s *Service) UnsubscribeMemo(ctx context.Context, memoID uuid.UUID, requester *model.AppIDToken) error {
	if requester == nil {
		return pkgerr.Known{
			Code:      pkgerr.CodeUnauthenticated,
			ClientMsg: "must be signed-in for subscription",
		}
	}

	err := s.transactionManager.WithTx(ctx, func(ctx context.Context) error {
		memo, err := s.memoRepository.FindByID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding memo: %w", err)
		}
		if memo.OwnerID == requester.UserID {
			return pkgerr.Known{
				Code:      pkgerr.CodeBadRequest,
				ClientMsg: "cannot unsubscribe memo of your own",
			}
		}

		subscribers, err := s.userRepository.FindAllBySubscribingMemoID(ctx, memoID)
		if err != nil {
			return fmt.Errorf("finding existing subscribers: %w", err)
		}

		if _, ok := lo.Find(subscribers, func(u *ent.User) bool { return u.ID == requester.UserID }); !ok {
			return pkgerr.Known{
				Code:      pkgerr.CodeBadRequest,
				ClientMsg: "memo is not subscribed",
			}
		}

		if err := s.memoRepository.UnregisterSubscriber(ctx, memoID, requester.UserID); err != nil {
			return fmt.Errorf("unregistering subscriber: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("during transaction: %w", err)
	}

	return nil
}

func (s *Service) ensureTags(ctx context.Context, tagNames []string) ([]*ent.Tag, error) {
	if err := validateTags(tagNames); err != nil {
		return nil, fmt.Errorf("validating tags: %w", err)
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

func validateTags(tags []string) error {
	for _, tag := range tags {
		if utf8.RuneCountInString(tag) > 20 {
			return pkgerr.Known{
				Code:      pkgerr.CodeBadRequest,
				ClientMsg: "length of tag should be less or equal to 20",
			}
		}

		if strings.TrimSpace(tag) == "" {
			return pkgerr.Known{
				Code:      pkgerr.CodeBadRequest,
				ClientMsg: "tag should not be blank",
			}
		}
	}

	return nil
}

func sortDedupTags(tags []string) []string {
	tags = lo.Uniq(tags)
	slices.Sort(tags)
	return tags
}

func removeReservedTags(tags []string) []string {
	removed, _ := lo.Difference(tags, reservedTags)
	return removed
}
