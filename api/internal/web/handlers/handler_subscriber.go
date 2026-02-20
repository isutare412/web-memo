package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/model"
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

	si, ok := lo.Find(resp.Subscribers, func(si model.SubscriberInfo) bool { return si.User.ID == userID })
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			ClientMsg: "no such subscriber",
		})
		return
	}

	gen.RespondJSON(w, http.StatusOK, gen.Subscriber{
		ID:       si.User.ID,
		UserName: si.User.UserName,
		PhotoURL: si.User.PhotoURL,
		Approved: si.Approved,
	})
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

	sub, err := h.memoService.SubscribeMemo(ctx, memoID, passport.Token)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("subscribing memo: %w", err))
		return
	}

	gen.RespondJSON(w, http.StatusOK, gen.SubscribeResponse{
		Subscription: struct {
			Approved bool               `json:"approved"`
			MemoID   openapi_types.UUID `json:"memoId"`
			UserID   openapi_types.UUID `json:"userId"`
		}{
			UserID:   sub.UserID,
			MemoID:   sub.MemoID,
			Approved: sub.Approved,
		},
	})
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

// AuthorizeSubscription approves or rejects a subscription request.
func (h *Handler) AuthorizeSubscription(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath, userID gen.UserIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.AuthorizeSubscription")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	var req gen.AuthorizeSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}

	if err := h.memoService.AuthorizeSubscriber(ctx, memoID, userID, req.Approve, passport.Token); err != nil {
		gen.RespondError(w, r, fmt.Errorf("authorizing subscriber: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}
