package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/isutare412/web-memo/api/internal/core/port"
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
	ctx := r.Context()

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
	ctx := r.Context()

	passport, ok := extractPassport(ctx)
	if !ok {
		responsePassportError(w, r)
		return
	}

	newToken, err := h.authService.RefreshAppIDToken(ctx, passport.tokenString)
	if err != nil {
		responseError(w, r, fmt.Errorf("refreshing app id token: %w", err))
		return
	}

	http.SetCookie(w, newWebMemoCookie(newToken, h.cookieExpiration))
	responseStatusCode(w, http.StatusOK)
}

func (h *userHandler) signOutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieNameWebMemoToken,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	})
}
