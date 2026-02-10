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

// ListSubscribers lists all subscribers for a memo.
func (h *Handler) ListSubscribers(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.ListSubscribers")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	resp, err := h.memoService.ListSubscribers(ctx, memoID, passport.Token)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("listing subscribers: %w", err))
		return
	}

	if resp.MemoOwnerID != passport.Token.UserID {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "cannot list subscribers of other user's memo",
		})
		return
	}

	gen.RespondJSON(w, http.StatusOK, SubscribersToWeb(resp))
}

// GetSubscriber returns a specific subscriber for a memo.
func (h *Handler) GetSubscriber(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath, userID gen.UserIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.GetSubscriber")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}
	if passport.Token.UserID != userID {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "token and userID mismatch",
		})
		return
	}

	resp, err := h.memoService.ListSubscribers(ctx, memoID, passport.Token)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("listing subscribers: %w", err))
		return
	}

	userFound, ok := lo.Find(resp.Subscribers, func(u *ent.User) bool { return u.ID == userID })
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			ClientMsg: "no such subscriber",
		})
		return
	}

	gen.RespondJSON(w, http.StatusOK, gen.Subscriber{ID: userFound.ID})
}

// Subscribe subscribes a user to a memo.
func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath, userID gen.UserIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.Subscribe")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}
	if passport.Token.UserID != userID {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "token and userID mismatch",
		})
		return
	}

	if err := h.memoService.SubscribeMemo(ctx, memoID, passport.Token); err != nil {
		gen.RespondError(w, r, fmt.Errorf("subscribing memo: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}

// Unsubscribe unsubscribes a user from a memo.
func (h *Handler) Unsubscribe(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath, userID gen.UserIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.Unsubscribe")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}
	if passport.Token.UserID != userID {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "token and userID mismatch",
		})
		return
	}

	if err := h.memoService.UnsubscribeMemo(ctx, memoID, passport.Token); err != nil {
		gen.RespondError(w, r, fmt.Errorf("subscribing memo: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}
