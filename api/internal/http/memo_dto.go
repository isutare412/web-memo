package http

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

type memoScores struct {
	RRF      float32 `json:"rrf"`
	Semantic float32 `json:"semantic"`
	BM25     float32 `json:"bm25"`
}

type memo struct {
	ID           uuid.UUID         `json:"id"`
	OwnerID      uuid.UUID         `json:"ownerId"`
	CreateTime   time.Time         `json:"createTime"`
	UpdateTime   time.Time         `json:"updateTime"`
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	PublishState enum.PublishState `json:"publishState"`
	Version      int               `json:"version"`
	Tags         []string          `json:"tags"`
	Scores       *memoScores       `json:"scores"`
}

func (m *memo) fromMemo(memo *ent.Memo) {
	m.ID = memo.ID
	m.OwnerID = memo.OwnerID
	m.CreateTime = memo.CreateTime
	m.UpdateTime = memo.UpdateTime
	m.Title = memo.Title
	m.Content = memo.Content
	m.PublishState = memo.PublishState
	m.Version = memo.Version
	m.Tags = lo.Map(memo.Edges.Tags, func(t *ent.Tag, _ int) string {
		return t.Name
	})
}

func (m *memo) fromMemoSearchResult(result *model.MemoSearchResult) {
	m.fromMemo(result.Memo)
	m.Scores = &memoScores{
		RRF:      result.RRFScore,
		Semantic: result.SemanticScore,
		BM25:     result.BM25Score,
	}
}

type listMemosResponse struct {
	Page           *int    `json:"page"`
	PageSize       *int    `json:"pageSize"`
	LastPage       *int    `json:"lastPage"`
	TotalMemoCount *int    `json:"totalMemoCount"`
	Memos          []*memo `json:"memos"`
}

type createMemoRequest struct {
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content" validate:"max=12000"`
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
	Content string   `json:"content" validate:"max=12000"`
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
	PublishState *enum.PublishState `json:"publishState" validate:"required"`
}

type listSubscribersResponse struct {
	Subscribers []*subscriber `json:"subscribers"`
}

type subscriber struct {
	ID       uuid.UUID `json:"id"`
	UserName string    `json:"userName"`
	PhotoURL string    `json:"photoUrl"`
	Approved bool      `json:"approved"`
}

func (s *subscriber) fromSubscriberInfo(si model.SubscriberInfo) {
	s.ID = si.User.ID
	s.UserName = si.User.UserName
	s.PhotoURL = si.User.PhotoURL
	s.Approved = si.Approved
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
