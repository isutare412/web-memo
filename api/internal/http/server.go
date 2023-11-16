package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/isutare412/web-memo/api/internal/core/port"
)

type Server struct {
	server *http.Server
}

func NewServer(cfg Config, authService port.AuthService, memoService port.MemoService, pingers []port.Pinger) *Server {
	healthHandler := newHealthHandler(pingers)
	googleHandler := newGoogleHandler(cfg, authService)
	userHandler := newUserHandler(authService)
	memoHandler := newMemoHandler(memoService)
	tagHandler := newTagHandler(memoService)

	imi := newImmigration(authService)

	r := chi.NewRouter()
	r.Mount("/health", healthHandler.router())

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(
			withContextBag,
			middleware.RealIP,
			wrapResponseWriter,
			logRequests,
			recoverPanic,
		)

		r.Mount("/google", googleHandler.router())

		auth := r.With(imi.issuePassport)
		auth.Mount("/users", userHandler.router())
		auth.Mount("/memos", memoHandler.router())
		auth.Mount("/tags", tagHandler.router())
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
