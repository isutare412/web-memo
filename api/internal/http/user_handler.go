package http

import (
	"fmt"
	"net/http"

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
	return r
}

func (h *userHandler) getSelfUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	passport, ok := extractPassport(ctx)
	if !ok {
		responseError(w, r, fmt.Errorf("passport not found"))
		return
	}

	var resp user
	resp.fromAppIDToken(passport.token)
	responseJSON(w, &resp)
}
