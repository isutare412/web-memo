package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/isutare412/web-memo/api/internal/core/port"
)

type userHandler struct {
	authService port.AuthService
}

func newUserHandler(authService port.AuthService) *userHandler {
	return &userHandler{
		authService: authService,
	}
}

func (h *userHandler) router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/me", h.getSelfUser)
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

func (h *userHandler) signOutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieNameWebMemoToken,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	})
}
