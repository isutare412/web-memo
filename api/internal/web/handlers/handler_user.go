package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/tracing"
	"github.com/isutare412/web-memo/api/internal/web/auth"
	"github.com/isutare412/web-memo/api/internal/web/gen"
	"github.com/isutare412/web-memo/api/internal/web/middleware"
)

// GetCurrentUser returns the currently authenticated user's information.
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.GetCurrentUser")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	gen.RespondJSON(w, http.StatusOK, UserToWeb(passport.Token))
}

// RefreshToken refreshes the authentication token and sets a new cookie.
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.RefreshToken")
	defer span.End()

	passport, ok := middleware.ExtractPassport(ctx)
	if !ok {
		gen.RespondError(w, r, pkgerr.Known{Code: pkgerr.CodeUnauthenticated, ClientMsg: "need token"})
		return
	}

	newToken, newTokenString, err := h.authService.RefreshAppIDToken(ctx, passport.TokenString)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("refreshing app id token: %w", err))
		return
	}

	http.SetCookie(w, auth.NewWebMemoCookie(newTokenString, h.cookieExpiration))
	gen.RespondJSON(w, http.StatusOK, UserToWeb(newToken))
}

// SignOut clears the authentication cookie.
func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	_, span := tracing.StartSpan(r.Context(), "web.handlers.SignOut")
	defer span.End()

	http.SetCookie(w, &http.Cookie{
		Name:    auth.CookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	})
}
