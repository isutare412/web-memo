package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/pkgerr"
)

const cookieNameWebMemoToken = "wmToken"

type contextKeyPassport struct{}

type passport struct {
	token *model.AppIDToken
}

type immigration struct {
	authService port.AuthService
}

func newImmigration(authService port.AuthService) *immigration {
	return &immigration{
		authService: authService,
	}
}

func (imi *immigration) issuePassport(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		passport, found, err := imi.issuePasspotFromHeader(w, r)
		if err != nil {
			responseError(w, r, fmt.Errorf("injecting passport from header: %w", err))
			return
		}

		if !found {
			passport, found, err = imi.issuePasspotFromCookie(w, r)
			switch {
			case err != nil:
				responseError(w, r, fmt.Errorf("injecting passport from cookie: %w", err))
				return
			case !found:
				responseError(w, r, pkgerr.Known{
					Code:      pkgerr.CodeUnauthenticated,
					ClientMsg: "need token",
				})
				return
			}
		}

		ctx := injectPassport(r.Context(), passport)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (imi *immigration) issuePasspotFromHeader(w http.ResponseWriter, r *http.Request) (*passport, bool, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil, false, nil
	}

	splitted := strings.SplitN(auth, " ", 2)
	if len(splitted) != 2 || splitted[0] != "Bearer" {
		return nil, false, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			ClientMsg: "authorization header must be Bearer format",
		}
	}

	token, err := imi.authService.VerifyAppIDTokenString(splitted[1])
	if err != nil {
		return nil, false, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("verifying app ID token: %w", err),
			ClientMsg: "invalid token",
		}
	}

	return &passport{
		token: token,
	}, true, nil
}

func (imi *immigration) issuePasspotFromCookie(w http.ResponseWriter, r *http.Request) (*passport, bool, error) {
	cookie, err := r.Cookie(cookieNameWebMemoToken)
	switch {
	case errors.Is(err, http.ErrNoCookie):
		return nil, false, nil
	case err != nil:
		return nil, false, fmt.Errorf("getting token cookie: %w", err)
	}

	token, err := imi.authService.VerifyAppIDTokenString(cookie.Value)
	if err != nil {
		return nil, false, pkgerr.Known{
			Code:      pkgerr.CodeBadRequest,
			Origin:    fmt.Errorf("verifying app ID token: %w", err),
			ClientMsg: "invalid token",
		}
	}

	return &passport{
		token: token,
	}, true, nil
}

func injectPassport(ctx context.Context, p *passport) context.Context {
	return context.WithValue(ctx, contextKeyPassport{}, p)
}

func extractPassport(ctx context.Context) (*passport, bool) {
	p, ok := ctx.Value(contextKeyPassport{}).(*passport)
	return p, ok
}
