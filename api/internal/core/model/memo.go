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

type ListSubscribersResponse struct {
	MemoOwnerID uuid.UUID
	Subscribers []*ent.User
}

type ListCollaboratorsResponse struct {
	MemoOwnerID   uuid.UUID
	Collaborators []*ent.User
}
