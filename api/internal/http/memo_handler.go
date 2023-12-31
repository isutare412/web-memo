package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/validate"
)

type memoHandler struct {
	memoService port.MemoService
}

func newMemoHandler(memoService port.MemoService) *memoHandler {
	return &memoHandler{
		memoService: memoService,
	}
}

func (h *memoHandler) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{memoID}", h.getMemo)
	r.Get("/", h.listMemos)
	r.Post("/", h.createMemo)
	r.Put("/{memoID}", h.replaceMemo)
	r.Delete("/{memoID}", h.deleteMemo)
	r.Get("/{memoID}/tags", h.listMemoTags)
	r.Put("/{memoID}/tags", h.replaceMemoTags)

	return r
}

func (h *memoHandler) getMemo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responseError(w, r, fmt.Errorf("passport not found"))
		return
	}

	memoFound, err := h.memoService.GetMemo(ctx, memoID, passport.token)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo: %w", err))
		return
	}

	var resp memo
	resp.fromMemo(memoFound)
	responseJSON(w, &resp)
}

func (h *memoHandler) listMemos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	passport, ok := extractPassport(ctx)
	if !ok {
		responseError(w, r, fmt.Errorf("passport not found"))
		return
	}

	page, pageSize, err := getPageQuery(r.URL.Query())
	if err != nil {
		responseError(w, r, fmt.Errorf("getting page query: %w", err))
		return
	}

	tags := getTagsQuery(r.URL.Query())

	queryOption := model.QueryOption{
		PageOffset: page,
		PageSize:   pageSize,
		Direction:  model.SortDirectionDesc,
	}

	memos, totalCount, err := h.memoService.ListMemos(ctx, passport.token.UserID, tags, &queryOption)
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
		Page:           page,
		PageSize:       pageSize,
		LastPage:       lastPage,
		TotalMemoCount: totalCount,
		Memos: lo.Map(memos, func(m *ent.Memo, _ int) *memo {
			var dto memo
			dto.fromMemo(m)
			return &dto
		}),
	}
	responseJSON(w, &resp)
}

func (h *memoHandler) createMemo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	passport, ok := extractPassport(ctx)
	if !ok {
		responseError(w, r, fmt.Errorf("passport not found"))
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
	ctx := r.Context()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responseError(w, r, fmt.Errorf("passport not found"))
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

	memoUpdated, err := h.memoService.UpdateMemo(ctx, memoToUpdate, req.Tags, passport.token)
	if err != nil {
		responseError(w, r, fmt.Errorf("updating memo: %w", err))
		return
	}

	var resp memo
	resp.fromMemo(memoUpdated)
	responseJSON(w, &resp)
}

func (h *memoHandler) deleteMemo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responseError(w, r, fmt.Errorf("passport not found"))
		return
	}

	if err := h.memoService.DeleteMemo(ctx, memoID, passport.token); err != nil {
		responseError(w, r, fmt.Errorf("deleting memo: %w", err))
		return
	}

	responseStatusCode(w, http.StatusOK)
}

func (h *memoHandler) replaceMemoTags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responseError(w, r, fmt.Errorf("passport not found"))
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
	ctx := r.Context()

	memoID, err := getMemoID(r)
	if err != nil {
		responseError(w, r, fmt.Errorf("getting memo ID: %w", err))
		return
	}

	passport, ok := extractPassport(ctx)
	if !ok {
		responseError(w, r, fmt.Errorf("passport not found"))
		return
	}

	tagsFound, err := h.memoService.ListTags(ctx, memoID, passport.token)
	if err != nil {
		responseError(w, r, fmt.Errorf("listing tags: %w", err))
		return
	}

	resp := lo.Map(tagsFound, func(tag *ent.Tag, _ int) string { return tag.Name })
	responseJSON(w, &resp)
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
