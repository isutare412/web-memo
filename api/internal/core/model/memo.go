package model

import (
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
)

type PaginationParams struct {
	PageOffset int
	PageSize   int
}

func (p *PaginationParams) GetOrDefault() (page, pageSize int) {
	page = lo.Ternary(p.PageOffset > 0, p.PageOffset, 1)
	pageSize = lo.Ternary(p.PageSize > 0, p.PageSize, 100)
	return page, pageSize
}

type MemoSortParams struct {
	Key   enum.MemoSortKey
	Order enum.SortOrder
}

type MemoSearchResult struct {
	Memo          *ent.Memo
	RRFScore      float32
	SemanticScore float32
	BM25Score     float32
}

type SubscriberInfo struct {
	User     *ent.User
	Approved bool
}

type ListSubscribersResponse struct {
	MemoOwnerID uuid.UUID
	Subscribers []SubscriberInfo
}

type ListCollaboratorsResponse struct {
	MemoOwnerID   uuid.UUID
	Collaborators []*ent.User
}

type MemoViewerContext struct {
	IsCollaborator bool
	IsApproved     bool
	Subscription   *ViewerSubscription
}

type ViewerSubscription struct {
	IsApproved bool
}

type GetMemoDetailResponse struct {
	Memo          *ent.Memo
	ViewerContext *MemoViewerContext
}
