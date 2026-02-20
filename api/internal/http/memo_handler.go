package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/tracing"
	"github.com/isutare412/web-memo/api/internal/validate"
)

type memoHandler struct {
	memoService      port.MemoService
	embeddingEnabled bool
}

func newMemoHandler(memoService port.MemoService, embeddingEnabled bool) *memoHandler {
	return &memoHandler{
		memoService:      memoService,
		embeddingEnabled: embeddingEnabled,
	}
}

func (h *memoHandler) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{memoID}", h.getMemo)
	r.Get("/", h.listMemos)
	r.Post("/", h.createMemo)
	r.Put("/{memoID}", h.replaceMemo)
	r.Post("/{memoID}/publish", h.publishMemo)
	r.Delete("/{memoID}", h.deleteMemo)
	r.Get("/{memoID}/tags", h.listMemoTags)
	r.Put("/{memoID}/tags", h.replaceMemoTags)
	r.Get("/{memoID}/subscribers/{userID}", h.getSubscriber)
	r.Get("/{memoID}/subscribers", h.listSubscribers)
	r.Put("/{memoID}/subscribers/{userID}", h.subscribeMemo)
	r.Delete("/{memoID}/subscribers/{userID}", h.unsubscribeMemo)
	r.Get("/{memoID}/collaborators/{userID}", h.getCollaborator)
	r.Get("/{memoID}/collaborators", h.listCollaborators)
	r.Put("/{memoID}/collaborators/{userID}", h.requestCollaboration)
	r.Post("/{memoID}/collaborators/{userID}/authorize", h.authorizeCollaboration)
	r.Delete("/{memoID}/collaborators/{userID}", h.cancelCollaboration)

	return r
}

func (h *memoHandler) getMemo(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.getMemo")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	var token *model.AppIDToken
	if passport, ok := extractPassport(ctx); ok {
		token = passport.token
	}

	memoFound, err := h.memoService.GetMemo(ctx, memoID, token)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo: %w", err))
		return
	}

	var resp memo
	resp.fromMemo(memoFound)
	responseJSON(w, &resp)
}

func (h *memoHandler) listMemos(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.listMemos")
	defer span.End()

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	query := r.URL.Query().Get("q")
	if query != "" {
		h.searchMemos(w, r, passport, query)
		return
	}

	page, pageSize, err := getPageQuery(r.URL.Query())
	if err != nil {
		responseError(w, r, fmt.Errorf("getting page query: %w", err))
		return
	}

	tags := getTagsQuery(r.URL.Query())
	sortKey := getMemoSortQuery(r.URL.Query())

	sortParams := model.MemoSortParams{
		Key:   sortKey,
		Order: enum.SortOrderDesc,
	}

	pageParams := model.PaginationParams{
		PageOffset: page,
		PageSize:   pageSize,
	}

	memos, totalCount, err := h.memoService.ListMemos(ctx, passport.token.UserID, tags, sortParams, pageParams)
	if err != nil {
		responseError(w, r, fmt.Errorf("listing memos: %w", err))
		return
	}

	lastPage := lo.Ternary(totalCount == 0, 1, (totalCount+pageSize-1)/pageSize)
	if page > lastPage {
		responseError(w, r, pkgerr.Known{Code: pkgerr.CodeNotFound, ClientMsg: "page not found"})
		return
	}

	resp := listMemosResponse{
		Page:           &page,
		PageSize:       &pageSize,
		LastPage:       &lastPage,
		TotalMemoCount: &totalCount,
		Memos: lo.Map(memos, func(m *ent.Memo, _ int) *memo {
			var dto memo
			dto.fromMemo(m)
			return &dto
		}),
	}
	responseJSON(w, &resp)
}

func (h *memoHandler) searchMemos(w http.ResponseWriter, r *http.Request, passport *passport, query string) {
	ctx := r.Context()

	if !h.embeddingEnabled {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "search feature is not enabled",
		})
		return
	}

	results, err := h.memoService.SearchMemos(ctx, passport.token.UserID, query)
	if err != nil {
		responseError(w, r, fmt.Errorf("searching memos: %w", err))
		return
	}

	memos := make([]*memo, 0, len(results))
	for _, result := range results {
		var dto memo
		dto.fromMemoSearchResult(result)
		memos = append(memos, &dto)
	}

	resp := listMemosResponse{
		Memos: memos,
	}
	responseJSON(w, &resp)
}

func (h *memoHandler) createMemo(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.createMemo")
	defer span.End()

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	var req createMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}
	if err := validate.Struct(&req); err != nil {
		responseError(w, r, fmt.Errorf("validating request body: %w", err))
		return
	}
	if err := req.validate(); err != nil {
		responseError(w, r, fmt.Errorf("validating request body: %w", err))
		return
	}

	memoCreated, err := h.memoService.CreateMemo(ctx, req.toMemo(), req.Tags, passport.token.UserID)
	if err != nil {
		responseError(w, r, fmt.Errorf("creating memo: %w", err))
		return
	}

	var resp memo
	resp.fromMemo(memoCreated)
	responseJSON(w, &resp)
}

func (h *memoHandler) replaceMemo(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.replaceMemo")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	isPinUpdateTime := getIsPinUpdateTime(r.URL.Query())

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	var req replaceMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}
	if err := validate.Struct(&req); err != nil {
		responseError(w, r, fmt.Errorf("validating request body: %w", err))
		return
	}
	if err := req.validate(); err != nil {
		responseError(w, r, fmt.Errorf("validating request body: %w", err))
		return
	}

	memoToUpdate := req.toMemo()
	memoToUpdate.ID = memoID

	memoUpdated, err := h.memoService.UpdateMemo(ctx, memoToUpdate, req.Tags, passport.token, isPinUpdateTime)
	if err != nil {
		responseError(w, r, fmt.Errorf("updating memo: %w", err))
		return
	}

	var resp memo
	resp.fromMemo(memoUpdated)
	responseJSON(w, &resp)
}

func (h *memoHandler) publishMemo(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.publishMemo")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	var req publishMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}
	if err := validate.Struct(&req); err != nil {
		responseError(w, r, fmt.Errorf("validating request body: %w", err))
		return
	}

	memoUpdated, err := h.memoService.UpdateMemoPublishState(ctx, memoID, *req.PublishState, passport.token)
	if err != nil {
		responseError(w, r, fmt.Errorf("updating memo published state: %w", err))
		return
	}

	var resp memo
	resp.fromMemo(memoUpdated)
	responseJSON(w, &resp)
}

func (h *memoHandler) deleteMemo(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.deleteMemo")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	if err := h.memoService.DeleteMemo(ctx, memoID, passport.token); err != nil {
		responseError(w, r, fmt.Errorf("deleting memo: %w", err))
		return
	}

	responseStatusCode(w, http.StatusOK)
}

func (h *memoHandler) replaceMemoTags(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.replaceMemoTags")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	var req []string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}

	tags, err := h.memoService.ReplaceTags(ctx, memoID, req, passport.token)
	if err != nil {
		responseError(w, r, fmt.Errorf("replacing tags: %w", err))
		return
	}

	resp := lo.Map(tags, func(tag *ent.Tag, _ int) string { return tag.Name })
	responseJSON(w, &resp)
}

func (h *memoHandler) listMemoTags(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.listMemoTags")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	var token *model.AppIDToken
	if passport, ok := extractPassport(ctx); ok {
		token = passport.token
	}

	tagsFound, err := h.memoService.ListTags(ctx, memoID, token)
	if err != nil {
		responseError(w, r, fmt.Errorf("listing tags: %w", err))
		return
	}

	resp := lo.Map(tagsFound, func(tag *ent.Tag, _ int) string { return tag.Name })
	responseJSON(w, &resp)
}

func (h *memoHandler) getSubscriber(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.getSubscriber")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting user ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}
	if passport.token.UserID != userID {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "token and userID mismatch",
		})
		return
	}

	resp, err := h.memoService.ListSubscribers(ctx, memoID, passport.token)
	if err != nil {
		responseError(w, r, fmt.Errorf("listing subscribers: %w", err))
		return
	}

	si, ok := lo.Find(resp.Subscribers, func(si model.SubscriberInfo) bool { return si.User.ID == userID })
	if !ok {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			ClientMsg: "no such subscriber",
		})
		return
	}

	var s subscriber
	s.fromSubscriberInfo(si)
	responseJSON(w, &s)
}

func (h *memoHandler) listSubscribers(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.listSubscribers")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	resp, err := h.memoService.ListSubscribers(ctx, memoID, passport.token)
	if err != nil {
		responseError(w, r, fmt.Errorf("listing subscribers: %w", err))
		return
	}

	if resp.MemoOwnerID != passport.token.UserID {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "cannot list subscribers of other user's memo",
		})
		return
	}

	responseJSON(w, &listSubscribersResponse{
		Subscribers: lo.Map(resp.Subscribers, func(si model.SubscriberInfo, _ int) *subscriber {
			var s subscriber
			s.fromSubscriberInfo(si)
			return &s
		}),
	})
}

func (h *memoHandler) subscribeMemo(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.subscribeMemo")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting user ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}
	if passport.token.UserID != userID {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "token and userID mismatch",
		})
		return
	}

	if _, err := h.memoService.SubscribeMemo(ctx, memoID, passport.token); err != nil {
		responseError(w, r, fmt.Errorf("subscribing memo: %w", err))
		return
	}

	responseStatusCode(w, http.StatusOK)
}

func (h *memoHandler) unsubscribeMemo(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.unsubscribeMemo")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting user ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}
	if passport.token.UserID != userID {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "token and userID mismatch",
		})
		return
	}

	if err := h.memoService.UnsubscribeMemo(ctx, memoID, passport.token); err != nil {
		responseError(w, r, fmt.Errorf("subscribing memo: %w", err))
		return
	}

	responseStatusCode(w, http.StatusOK)
}

func (h *memoHandler) getCollaborator(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.getCollaborator")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting user ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}
	if passport.token.UserID != userID {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "token and userID mismatch",
		})
		return
	}

	resp, err := h.memoService.ListCollaborators(ctx, memoID, passport.token)
	if err != nil {
		responseError(w, r, fmt.Errorf("listing subscribers: %w", err))
		return
	}

	user, ok := lo.Find(resp.Collaborators, func(u *ent.User) bool { return u.ID == userID })
	if !ok {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			ClientMsg: "no such collaborator",
		})
		return
	}

	collabo, ok := lo.Find(user.Edges.Collaborations, func(c *ent.Collaboration) bool { return c.MemoID == memoID })
	if !ok {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			ClientMsg: "no such collaboration",
		})
		return
	}

	responseJSON(w, &collaborator{
		ID:         user.ID,
		UserName:   user.UserName,
		PhotoURL:   user.PhotoURL,
		IsApproved: collabo.Approved,
	})
}

func (h *memoHandler) listCollaborators(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.listCollaborators")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	resp, err := h.memoService.ListCollaborators(ctx, memoID, passport.token)
	if err != nil {
		responseError(w, r, fmt.Errorf("listing subscribers: %w", err))
		return
	}

	if resp.MemoOwnerID != passport.token.UserID {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "cannot list collaborators of other user's memo",
		})
		return
	}

	collaborators := make([]*collaborator, 0, len(resp.Collaborators))
	for _, user := range resp.Collaborators {
		for _, collabo := range user.Edges.Collaborations {
			if collabo.MemoID != memoID {
				continue
			}

			var c collaborator
			c.fromModels(collabo, user)
			collaborators = append(collaborators, &c)
		}
	}

	responseJSON(w, &listCollaboratorsResponse{
		Collaborators: collaborators,
	})
}

func (h *memoHandler) requestCollaboration(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.requestCollaboration")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting user ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}
	if passport.token.UserID != userID {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "token and userID mismatch",
		})
		return
	}

	if err := h.memoService.RegisterCollaborator(ctx, memoID, passport.token); err != nil {
		responseError(w, r, fmt.Errorf("registering collaborator: %w", err))
		return
	}

	responseStatusCode(w, http.StatusOK)
}

func (h *memoHandler) authorizeCollaboration(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.authorizeCollaboration")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting user ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	var req authorizeCollaborationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}
	if err := validate.Struct(&req); err != nil {
		responseError(w, r, fmt.Errorf("validating request body: %w", err))
		return
	}

	if err := h.memoService.AuthorizeCollaborator(ctx, memoID, userID, *req.Approve, passport.token); err != nil {
		responseError(w, r, fmt.Errorf("authorizing collaborator: %w", err))
		return
	}

	responseStatusCode(w, http.StatusOK)
}

func (h *memoHandler) cancelCollaboration(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.memoHandler.cancelCollaboration")
	defer span.End()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting user ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}
	if passport.token.UserID != userID {
		responseError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "token and userID mismatch",
		})
		return
	}

	if err := h.memoService.DeleteCollaborator(ctx, memoID, userID, passport.token); err != nil {
		responseError(w, r, fmt.Errorf("authorizing collaborator: %w", err))
		return
	}

	responseStatusCode(w, http.StatusOK)
}

func getMemoID(r *http.Request) (uuid.UUID, error) {
	memoIDString := chi.URLParam(r, "memoID")
	if memoIDString == "" {
		return uuid.UUID{}, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "need memoID",
		}
	}

	memoID, err := uuid.Parse(memoIDString)
	if err != nil {
		return uuid.UUID{}, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "format of memo ID is invalid",
		}
	}

	return memoID, nil
}

func getUserID(r *http.Request) (uuid.UUID, error) {
	userIDString := chi.URLParam(r, "userID")
	if userIDString == "" {
		return uuid.UUID{}, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "need userID",
		}
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.UUID{}, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "format of user ID is invalid",
		}
	}

	return userID, nil
}

func getPageQuery(q url.Values) (page, pageSize int, err error) {
	pageStr := q.Get("page")
	pageSizeStr := q.Get("pageSize")

	page = 1
	pageSize = 10

	if pageStr != "" {
		n, err := strconv.Atoi(pageStr)
		if err != nil || n <= 0 {
			return 0, 0, pkgerr.Known{Code: pkgerr.CodeBadRequest, ClientMsg: "query page should be a positive number"}
		}
		page = n
	}

	if pageSizeStr != "" {
		n, err := strconv.Atoi(pageSizeStr)
		if err != nil || n <= 0 {
			return 0, 0, pkgerr.Known{Code: pkgerr.CodeBadRequest, ClientMsg: "query pageSize should be a positive number"}
		} else if n > 100 {
			return 0, 0, pkgerr.Known{Code: pkgerr.CodeBadRequest, ClientMsg: "query pageSize is too large"}
		}
		pageSize = n
	}

	return page, pageSize, nil
}

func getTagsQuery(q url.Values) []string {
	return q["tag"]
}

func getMemoSortQuery(q url.Values) enum.MemoSortKey {
	return enum.MemoSortKey(q.Get("sort")).GetOrDefault()
}

func getIsPinUpdateTime(q url.Values) bool {
	return strings.EqualFold(q.Get("pinUpdateTime"), "true")
}
