package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/isutare412/web-memo/api/internal/core/port"
)

type Server struct {
	server *http.Server
}

func NewServer(cfg Config, authService port.AuthService) *Server {
	ggHandler := newGoogleHandler(cfg, authService)
	usrHandler := newUserHandler(authService)

	imi := newImmigration(authService)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/google", ggHandler.router())

		auth := r.With(imi.issuePassport)
		auth.Mount("/users", usrHandler.router())
	})

	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
			Handler: r,
		},
	}
}

func (s *Server) Run() <-chan error {
	runtimeErrs := make(chan error, 1)
	go func() {
		slog.Info("run http server", "address", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			runtimeErrs <- fmt.Errorf("listen and serving HTTP: %w", err)
			return
		}
	}()

	return runtimeErrs
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
