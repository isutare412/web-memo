package web

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/metric"
	"github.com/isutare412/web-memo/api/internal/web/auth"
	"github.com/isutare412/web-memo/api/internal/web/gen"
	"github.com/isutare412/web-memo/api/internal/web/handlers"
	"github.com/isutare412/web-memo/api/internal/web/middleware"
)

//go:embed openapi.yaml openapi.html
var staticContents embed.FS

type Server struct {
	cfg    Config
	server *http.Server
}

func NewServer(
	cfg Config,
	pingers []port.Pinger,
	authService port.AuthService,
	memoService port.MemoService,
	imageService port.ImageService,
) *Server {
	handler := handlers.NewHandler(
		pingers, authService, memoService, imageService,
		cfg.CookieTokenExpiration, cfg.EmbeddingEnabled,
	)

	authenticator := auth.NewAuthenticator(authService)

	baseMiddlewares := []mux.MiddlewareFunc{
		middleware.WithLogAttrContext,
		middleware.WithTrace,
		middleware.WithContextBag,
		middleware.WithRequestID,
		middleware.ProxyHeaders,
		middleware.WithResponseRecord,
		middleware.AccessLog,
		middleware.ObserveMetrics,
		middleware.RecoverPanic,
	}

	apiMiddlewares := append([]mux.MiddlewareFunc(nil), baseMiddlewares...)
	apiMiddlewares = append(apiMiddlewares,
		authenticator.Authenticate,
		middleware.WithOpenAPIValidator(),
	)

	r := mux.NewRouter()

	// Health checks (base middleware only)
	healthRouter := r.PathPrefix("/health").Subrouter()
	healthRouter.Use(baseMiddlewares...)
	healthRouter.HandleFunc("/livez", handler.Liveness).Methods("GET")
	healthRouter.HandleFunc("/readyz", handler.Readiness).Methods("GET")

	// OpenAPI docs (base middleware only, no auth)
	if cfg.ShowOpenAPIDocs {
		docsRouter := r.PathPrefix("/docs").Subrouter()
		docsRouter.Use(baseMiddlewares...)

		docsRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/docs/openapi.html", http.StatusMovedPermanently)
		}).Methods("GET")

		docsRouter.PathPrefix("/").
			Handler(http.StripPrefix("/docs/",
				http.FileServer(http.FS(staticContents)))).
			Methods("GET")
	}

	// Metrics (no middleware)
	r.Handle("/metrics", promhttp.HandlerFor(metric.Gatherer(), promhttp.HandlerOpts{})).
		Methods("GET")

	// API routes (all middleware)
	apiRouter := r.PathPrefix("/").Subrouter()
	apiRouter.Use(apiMiddlewares...)
	gen.HandlerWithOptions(handler, gen.GorillaServerOptions{
		BaseRouter:       apiRouter,
		ErrorHandlerFunc: gen.RespondError,
	})

	return &Server{
		cfg: cfg,
		server: &http.Server{
			Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
			Handler: r,
		},
	}
}

func (s *Server) Run() <-chan error {
	runtimeErrs := make(chan error, 1)
	go func() {
		defer close(runtimeErrs)

		slog.Info("Starting web server", "port", s.cfg.Port)
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
