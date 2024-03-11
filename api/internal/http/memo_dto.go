package http

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

type memo struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"ownerId"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	IsPublished bool      `json:"isPublished"`
	Version     int       `json:"version"`
	Tags        []string  `json:"tags"`
}

func (m *memo) fromMemo(memo *ent.Memo) {
	m.ID = memo.ID
	m.OwnerID = memo.OwnerID
	m.CreateTime = memo.CreateTime
	m.UpdateTime = memo.UpdateTime
	m.Title = memo.Title
	m.Content = memo.Content
	m.IsPublished = memo.IsPublished
	m.Version = memo.Version
	m.Tags = lo.Map(memo.Edges.Tags, func(t *ent.Tag, _ int) string {
		return t.Name
	})
}

type listMemosResponse struct {
	Page           int     `json:"page"`
	PageSize       int     `json:"pageSize"`
	LastPage       int     `json:"lastPage"`
	TotalMemoCount int     `json:"totalMemoCount"`
	Memos          []*memo `json:"memos"`
}

type createMemoRequest struct {
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (r *createMemoRequest) validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "title should not be blank string",
		}
	}
	return nil
}

func (r *createMemoRequest) toMemo() *ent.Memo {
	return &ent.Memo{
		Title:   r.Title,
		Content: r.Content,
	}
}

type replaceMemoRequest struct {
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Version *int     `json:"version" validate:"required"`
}

func (r *replaceMemoRequest) validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "title should not be blank string",
		}
	}
	return nil
}

func (r *replaceMemoRequest) toMemo() *ent.Memo {
	return &ent.Memo{
		Title:   r.Title,
		Content: r.Content,
		Version: *r.Version,
	}
}

type publishMemoRequest struct {
	Publish *bool `json:"publish" validate:"required"`
}

type listSubscribersResponse struct {
	Subscribers []*subscriber `json:"subscribers"`
}

type subscriber struct {
	ID uuid.UUID `json:"id"`
}

func (s *subscriber) fromUser(u *ent.User) {
	s.ID = u.ID
}

type listCollaboratorsResponse struct {
	Collaborators []*collaborator `json:"collaborators"`
}

type collaborator struct {
	ID         uuid.UUID `json:"id"`
	UserName   string    `json:"userName"`
	PhotoURL   string    `json:"photoUrl"`
	IsApproved bool      `json:"isApproved"`
}

func (c *collaborator) fromModels(collabo *ent.Collaboration, user *ent.User) {
	c.ID = collabo.UserID
	c.UserName = user.UserName
	c.PhotoURL = user.PhotoURL
	c.IsApproved = collabo.Approved
}

type authorizeCollaborationRequest struct {
	Approve *bool `json:"approve" validate:"required"`
}
