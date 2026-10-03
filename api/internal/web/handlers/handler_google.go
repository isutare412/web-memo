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
func (h *Handler) StartGoogleSignIn(w http.ResponseWriter, r *http.Request, params gen.StartGoogleSignInParams) {
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

// FinishGoogleSignIn completes the Google OAuth2 sign-in flow and redirects
// the user. Web logins get the authentication cookie; CLI logins get the token
// in the loopback callback URL instead.
func (h *Handler) FinishGoogleSignIn(w http.ResponseWriter, r *http.Request, params gen.FinishGoogleSignInParams) {
	ctx, span := tracing.StartSpan(r.Context(), "web.handlers.FinishGoogleSignIn")
	defer span.End()

	result, err := h.authService.FinishGoogleSignIn(ctx, r)
	if err != nil {
		gen.RespondError(w, r, fmt.Errorf("finishing google sign-in: %w", err))
		return
	}

	if result.SetCookie {
		slog.Info("finished google sign-in", "redirectURL", result.RedirectURL)
		http.SetCookie(w, auth.NewWebMemoCookie(result.AppToken, h.cookieExpiration))
	} else {
		// The redirect URL of a CLI login carries the app token.
		slog.Info("finished google sign-in for cli")
	}
	http.Redirect(w, r, result.RedirectURL, http.StatusFound)
}
