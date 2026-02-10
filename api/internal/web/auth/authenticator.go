package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/web/gen"
	"github.com/isutare412/web-memo/api/internal/web/middleware"
)

// Authenticator extracts and verifies authentication tokens from HTTP requests.
type Authenticator struct {
	authService port.AuthService
}

// NewAuthenticator creates a new Authenticator with the given AuthService.
func NewAuthenticator(authService port.AuthService) *Authenticator {
	return &Authenticator{authService: authService}
}

// Authenticate is an HTTP middleware that extracts authentication tokens from
// the Authorization header or cookie, verifies them, and stores the result in
// the request context's ContextBag.
func (a *Authenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passport, found, err := a.authenticateFromHeader(r)
		if err != nil {
			gen.RespondError(w, r, fmt.Errorf("authenticating from header: %w", err))
			return
		}

		if !found {
			passport, found, err = a.authenticateFromCookie(r)
			if err != nil {
				gen.RespondError(w, r, fmt.Errorf("authenticating from cookie: %w", err))
				return
			}
		}

		if bag, ok := middleware.GetContextBag(r.Context()); ok && found {
			bag.Passport = passport
		}
		next.ServeHTTP(w, r)
	})
}

func (a *Authenticator) authenticateFromHeader(r *http.Request) (*middleware.Passport, bool, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil, false, nil
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, false, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "authorization header must be Bearer format",
		}
	}

	tokenString := parts[1]
	token, err := a.authService.VerifyAppIDToken(tokenString)
	if err != nil {
		return nil, false, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("verifying app ID token: %w", err),
			ClientMsg: "invalid token",
		}
	}

	return &middleware.Passport{
		TokenString: tokenString,
		Token:       token,
	}, true, nil
}

func (a *Authenticator) authenticateFromCookie(r *http.Request) (*middleware.Passport, bool, error) {
	cookie, err := r.Cookie(CookieName)
	switch {
	case errors.Is(err, http.ErrNoCookie):
		return nil, false, nil
	case err != nil:
		return nil, false, fmt.Errorf("getting token cookie: %w", err)
	}

	tokenString := cookie.Value
	token, err := a.authService.VerifyAppIDToken(tokenString)
	if err != nil {
		return nil, false, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("verifying app ID token: %w", err),
			ClientMsg: "invalid token",
		}
	}

	return &middleware.Passport{
		TokenString: tokenString,
		Token:       token,
	}, true, nil
}
