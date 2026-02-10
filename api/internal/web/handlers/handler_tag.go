package handlers

import (
	"fmt"
	"net/http"

	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/tracing"
	"github.com/isutare412/web-memo/api/internal/web/gen"
	"github.com/isutare412/web-memo/api/internal/web/middleware"
)

// SearchTags searches tags by keyword.
func (h *Handler) SearchTags(w http.ResponseWriter, r *http.Request, params gen.SearchTagsParams) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.SearchTags")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	var keyword string
	if params.Kw != nil {
		keyword = *params.Kw
	}

	tags, err := h.memoService.SearchTags(ctx, keyword, passport.Token)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("searching tags: %w", err))
		return
	}

	resp := lo.Map(tags, func(tag *ent.Tag, _ int) string { return tag.Name })
	gen.RespondJSON(w, http.StatusOK, resp)
}
