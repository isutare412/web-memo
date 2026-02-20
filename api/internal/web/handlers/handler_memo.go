package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/enum"
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/tracing"
	"github.com/isutare412/web-memo/api/internal/web/gen"
	"github.com/isutare412/web-memo/api/internal/web/middleware"
)

// ListMemos returns a paginated list of memos or search results.
func (h *Handler) ListMemos(w http.ResponseWriter, r *http.Request, params gen.ListMemosParams) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.ListMemos")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	if params.Q != nil && *params.Q != "" {
		h.searchMemos(w, r, passport, *params.Q)
		return
	}

	page := 1
	if params.Page != nil {
		page = *params.Page
	}

	pageSize := 10
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	tags := params.Tag

	var sortKey enum.MemoSortKey
	if params.Sort != nil {
		sortKey = *params.Sort
	}
	sortKey = sortKey.GetOrDefault()

	sortParams := model.MemoSortParams{
		Key:   sortKey,
		Order: enum.SortOrderDesc,
	}

	pageParams := model.PaginationParams{
		PageOffset: page,
		PageSize:   pageSize,
	}

	memos, totalCount, err := h.memoService.ListMemos(ctx, passport.Token.UserID, tags, sortParams, pageParams)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("listing memos: %w", err))
		return
	}

	lastPage := lo.Ternary(totalCount == 0, 1, (totalCount+pageSize-1)/pageSize)
	if page > lastPage {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeNotFound, ClientMsg: "page not found"})
		return
	}

	tagsByMemo := make([][]string, len(memos))
	for i, m := range memos {
		tagsByMemo[i] = lo.Map(m.Edges.Tags, func(t *ent.Tag, _ int) string { return t.Name })
	}

	resp := MemosToListResponse(memos, tagsByMemo, page, pageSize, lastPage, totalCount)
	gen.RespondJSON(w, http.StatusOK, resp)
}

func (h *Handler) searchMemos(w http.ResponseWriter, r *http.Request, passport *middleware.Passport, query string) {
	ctx := r.Context()

	if !h.embeddingEnabled {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "search feature is not enabled",
		})
		return
	}

	results, err := h.memoService.SearchMemos(ctx, passport.Token.UserID, query)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("searching memos: %w", err))
		return
	}

	resp := SearchResultsToListResponse(results)
	gen.RespondJSON(w, http.StatusOK, resp)
}

// CreateMemo creates a new memo.
func (h *Handler) CreateMemo(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.CreateMemo")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	var req gen.CreateMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "title should not be blank string",
		})
		return
	}

	memo := &ent.Memo{
		Title:   req.Title,
		Content: lo.FromPtrOr(req.Content, ""),
	}

	tags := req.Tags

	memoCreated, err := h.memoService.CreateMemo(ctx, memo, tags, passport.Token.UserID)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("creating memo: %w", err))
		return
	}

	memoTags := lo.Map(memoCreated.Edges.Tags, func(t *ent.Tag, _ int) string { return t.Name })
	gen.RespondJSON(w, http.StatusOK, MemoToWeb(memoCreated, memoTags))
}

// GetMemo returns a single memo by ID.
func (h *Handler) GetMemo(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.GetMemo")
	defer span.End()

	var token *model.AppIDToken
	if passport, ok := middleware.ExtractPassport(ctx); ok {
		token = passport.Token
	}

	resp, err := h.memoService.GetMemoDetail(ctx, memoID, token)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("getting memo detail: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, MemoDetailToWeb(resp))
}

// ReplaceMemo replaces a memo with updated data.
func (h *Handler) ReplaceMemo(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath, params gen.ReplaceMemoParams) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.ReplaceMemo")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	var req gen.ReplaceMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "title should not be blank string",
		})
		return
	}

	isPinUpdateTime := params.PinUpdateTime != nil && *params.PinUpdateTime

	memoToUpdate := &ent.Memo{
		ID:      memoID,
		Title:   req.Title,
		Content: lo.FromPtrOr(req.Content, ""),
		Version: req.Version,
	}

	tags := req.Tags

	memoUpdated, err := h.memoService.UpdateMemo(ctx, memoToUpdate, tags, passport.Token, isPinUpdateTime)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("updating memo: %w", err))
		return
	}

	memoTags := lo.Map(memoUpdated.Edges.Tags, func(t *ent.Tag, _ int) string { return t.Name })
	gen.RespondJSON(w, http.StatusOK, MemoToWeb(memoUpdated, memoTags))
}

// DeleteMemo deletes a memo by ID.
func (h *Handler) DeleteMemo(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.DeleteMemo")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	if err := h.memoService.DeleteMemo(ctx, memoID, passport.Token); err != nil {
		gen.RespondError(w, r, fmt.Errorf("deleting memo: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}

// PublishMemo publishes or unpublishes a memo.
func (h *Handler) PublishMemo(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.PublishMemo")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	var req gen.PublishMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}

	memoUpdated, err := h.memoService.UpdateMemoPublishState(ctx, memoID,
		enum.PublishState(req.PublishState), passport.Token)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("updating memo published state: %w", err))
		return
	}

	memoTags := lo.Map(memoUpdated.Edges.Tags, func(t *ent.Tag, _ int) string { return t.Name })
	gen.RespondJSON(w, http.StatusOK, MemoToWeb(memoUpdated, memoTags))
}

// GetMemoTags returns tags for a specific memo.
func (h *Handler) GetMemoTags(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.GetMemoTags")
	defer span.End()

	var token *model.AppIDToken
	if passport, ok := middleware.ExtractPassport(ctx); ok {
		token = passport.Token
	}

	tagsFound, err := h.memoService.ListTags(ctx, memoID, token)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("listing tags: %w", err))
		return
	}

	resp := lo.Map(tagsFound, func(tag *ent.Tag, _ int) string { return tag.Name })
	gen.RespondJSON(w, http.StatusOK, resp)
}

// ReplaceMemoTags replaces all tags of a memo.
func (h *Handler) ReplaceMemoTags(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.ReplaceMemoTags")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	var req []string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}

	tags, err := h.memoService.ReplaceTags(ctx, memoID, req, passport.Token)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("replacing tags: %w", err))
		return
	}

	resp := lo.Map(tags, func(tag *ent.Tag, _ int) string { return tag.Name })
	gen.RespondJSON(w, http.StatusOK, resp)
}
