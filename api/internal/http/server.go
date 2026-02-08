package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/metric"
)

type Server struct {
	server *http.Server
}

func NewServer(
	cfg Config,
	authService port.AuthService,
	memoService port.MemoService,
	pingers []port.Pinger,
	imageService port.ImageService,
) *Server {
	healthHandler := newHealthHandler(pingers)
	googleHandler := newGoogleHandler(cfg, authService)
	userHandler := newUserHandler(cfg, authService)
	memoHandler := newMemoHandler(memoService, cfg.EmbeddingEnabled)
	tagHandler := newTagHandler(memoService)
	imageHandler := newImageHandler(imageService)

	imi := newImmigration(authService)

	r := chi.NewRouter()

	r.Handle("/metrics", promhttp.HandlerFor(
		metric.Gatherer(),
		promhttp.HandlerOpts{},
	))

	r.Mount("/health", healthHandler.router())

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(
			withLogAttrContext,
			withTrace,
			withContextBag,
			middleware.RealIP,
			wrapResponseWriter,
			logRequests,
			observeMetrics,
			recoverPanic,
		)

		r.Mount("/google", googleHandler.router())

		auth := r.With(imi.issuePassport)
		auth.Mount("/users", userHandler.router())
		auth.Mount("/memos", memoHandler.router())
		auth.Mount("/tags", tagHandler.router())
		auth.Mount("/images", imageHandler.router())
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
