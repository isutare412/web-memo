package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/tracing"
)

type userHandler struct {
	authService      port.AuthService
	cookieExpiration time.Duration
}

func newUserHandler(cfg Config, authService port.AuthService) *userHandler {
	return &userHandler{
		authService:      authService,
		cookieExpiration: cfg.CookieTokenExpiration,
	}
}

func (h *userHandler) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/me", h.getSelfUser)
	r.Post("/refresh-token", h.refreshUserToken)
	r.Get("/sign-out", h.signOutUser)
	return r
}

func (h *userHandler) getSelfUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.userHandler.getSelfUser")
	defer span.End()

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	var resp user
	resp.fromAppIDToken(passport.token)
	responseJSON(w, &resp)
}

func (h *userHandler) refreshUserToken(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "http.userHandler.refreshUserToken")
	defer span.End()

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	newToken, newTokenString, err := h.authService.RefreshAppIDToken(ctx, passport.tokenString)
	if err != nil {
		responseError(w, r, fmt.Errorf("refreshing app id token: %w", err))
		return
	}

	http.SetCookie(w, newWebMemoCookie(newTokenString, h.cookieExpiration))
	var resp user
	resp.fromAppIDToken(newToken)
	responseJSON(w, &resp)
}

func (h *userHandler) signOutUser(w http.ResponseWriter, r *http.Request) {
	_, span := tracing.StartSpan(r.Context(), "http.userHandler.signOutUser")
	defer span.End()

	http.SetCookie(w, &http.Cookie{
		Name:    cookieNameWebMemoToken,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	})
}
