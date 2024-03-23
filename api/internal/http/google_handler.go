package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/isutare412/web-memo/api/internal/core/port"
)

type googleHandler struct {
	authService      port.AuthService
	cookieExpiration time.Duration
}

func newGoogleHandler(cfg Config, authService port.AuthService) *googleHandler {
	return &googleHandler{
		authService:      authService,
		cookieExpiration: cfg.CookieTokenExpiration,
	}
}

func (h *googleHandler) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/sign-in", h.googleSignIn)
	r.Get("/sign-in/finish", h.googleSignInFinish)
	return r
}

func (h *googleHandler) googleSignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	redirectURL, err := h.authService.StartGoogleSignIn(ctx, r)
	if err != nil {
		responseError(w, r, fmt.Errorf("starting google sign-in: %w", err))
		return
	}

	slog.Info("redirect user for google sign-in", "redirectUrl", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *googleHandler) googleSignInFinish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	redirectURL, appToken, err := h.authService.FinishGoogleSignIn(ctx, r)
	if err != nil {
		responseError(w, r, fmt.Errorf("finishing google sign-in: %w", err))
		return
	}

	slog.Info("finished google sign-in", "redirectURL", redirectURL)
	http.SetCookie(w, newWebMemoCookie(appToken, h.cookieExpiration))
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
