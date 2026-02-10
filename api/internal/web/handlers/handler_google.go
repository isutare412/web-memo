package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/isutare412/web-memo/api/internal/tracing"
	"github.com/isutare412/web-memo/api/internal/web/auth"
	"github.com/isutare412/web-memo/api/internal/web/gen"
)

// StartGoogleSignIn initiates the Google OAuth2 sign-in flow by redirecting
// the user to Google's authorization page.
func (h *Handler) StartGoogleSignIn(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.StartGoogleSignIn")
	defer span.End()

	redirectURL, err := h.authService.StartGoogleSignIn(ctx, r)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("starting google sign-in: %w", err))
		return
	}

	slog.Info("redirect user for google sign-in", "redirectUrl", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// FinishGoogleSignIn completes the Google OAuth2 sign-in flow, sets the
// authentication cookie, and redirects the user.
func (h *Handler) FinishGoogleSignIn(w http.ResponseWriter, r *http.Request, params gen.FinishGoogleSignInParams) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.FinishGoogleSignIn")
	defer span.End()

	redirectURL, appToken, err := h.authService.FinishGoogleSignIn(ctx, r)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("finishing google sign-in: %w", err))
		return
	}

	slog.Info("finished google sign-in", "redirectURL", redirectURL)
	http.SetCookie(w, auth.NewWebMemoCookie(appToken, h.cookieExpiration))
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
