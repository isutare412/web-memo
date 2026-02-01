package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/tracing"
)

type tagHandler struct {
	memoService port.MemoService
}

func newTagHandler(memoService port.MemoService) *tagHandler {
	return &tagHandler{
		memoService: memoService,
	}
}

func (h *tagHandler) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", h.listTags)

	return r
}

func (h *tagHandler) listTags(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.tagHandler.listTags")
	defer span.End()

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	keyword := r.URL.Query().Get("kw")

	tags, err := h.memoService.SearchTags(ctx, keyword, passport.token)
	if err != nil {
		responseError(w, r, fmt.Errorf("searching tags: %w", err))
		return
	}

	resp := lo.Map(tags, func(tag *ent.Tag, _ int) string { return tag.Name })
	responseJSON(w, &resp)
}
