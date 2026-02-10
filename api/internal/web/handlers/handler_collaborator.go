package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/tracing"
	"github.com/isutare412/web-memo/api/internal/web/gen"
	"github.com/isutare412/web-memo/api/internal/web/middleware"
)

// ListCollaborators lists all collaborators for a memo.
func (h *Handler) ListCollaborators(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.ListCollaborators")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	resp, err := h.memoService.ListCollaborators(ctx, memoID, passport.Token)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("listing collaborators: %w", err))
		return
	}

	if resp.MemoOwnerID != passport.Token.UserID {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodePermissionDenied,
			ClientMsg: "cannot list collaborators of other user's memo",
		})
		return
	}

	gen.RespondJSON(w, http.StatusOK, CollaboratorsToWeb(resp, memoID))
}

// GetCollaborator returns a specific collaborator for a memo.
func (h *Handler) GetCollaborator(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath, userID gen.UserIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.GetCollaborator")
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

	resp, err := h.memoService.ListCollaborators(ctx, memoID, passport.Token)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("listing collaborators: %w", err))
		return
	}

	user, ok := lo.Find(resp.Collaborators, func(u *ent.User) bool { return u.ID == userID })
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			ClientMsg: "no such collaborator",
		})
		return
	}

	collabo, ok := lo.Find(user.Edges.Collaborations, func(c *ent.Collaboration) bool { return c.MemoID == memoID })
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			ClientMsg: "no such collaboration",
		})
		return
	}

	gen.RespondJSON(w, http.StatusOK, gen.Collaborator{
		ID:         user.ID,
		UserName:   user.UserName,
		PhotoURL:   user.PhotoURL,
		IsApproved: collabo.Approved,
	})
}

// RequestCollaboration requests collaboration on a memo.
func (h *Handler) RequestCollaboration(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath, userID gen.UserIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.RequestCollaboration")
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

	if err := h.memoService.RegisterCollaborator(ctx, memoID, passport.Token); err != nil {
		gen.RespondError(w, r, fmt.Errorf("registering collaborator: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}

// AuthorizeCollaboration approves or rejects a collaboration request.
func (h *Handler) AuthorizeCollaboration(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath, userID gen.UserIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.AuthorizeCollaboration")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	var req gen.AuthorizeCollaborationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		gen.RespondError(w, r, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("decoding request body: %w", err),
			ClientMsg: "invalid request body",
		})
		return
	}

	if err := h.memoService.AuthorizeCollaborator(ctx, memoID, userID, req.Approve, passport.Token); err != nil {
		gen.RespondError(w, r, fmt.Errorf("authorizing collaborator: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}

// CancelCollaboration cancels a collaboration on a memo.
func (h *Handler) CancelCollaboration(w http.ResponseWriter, r *http.Request, memoID gen.MemoIDPath, userID gen.UserIDPath) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.CancelCollaboration")
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

	if err := h.memoService.DeleteCollaborator(ctx, memoID, userID, passport.Token); err != nil {
		gen.RespondError(w, r, fmt.Errorf("canceling collaboration: %w", err))
		return
	}

	gen.RespondNoContent(w, http.StatusOK)
}
