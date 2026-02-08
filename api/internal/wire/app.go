package wire

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isutare412/web-memo/api/internal/config"
	"github.com/isutare412/web-memo/api/internal/core/port"
	"github.com/isutare412/web-memo/api/internal/core/service/auth"
	"github.com/isutare412/web-memo/api/internal/core/service/image"
	"github.com/isutare412/web-memo/api/internal/core/service/memo"
	"github.com/isutare412/web-memo/api/internal/cron"
	"github.com/isutare412/web-memo/api/internal/embedding"
	"github.com/isutare412/web-memo/api/internal/google"
	"github.com/isutare412/web-memo/api/internal/http"
	"github.com/isutare412/web-memo/api/internal/imageer"
	"github.com/isutare412/web-memo/api/internal/jwt"
	"github.com/isutare412/web-memo/api/internal/postgres"
	"github.com/isutare412/web-memo/api/internal/redis"
)

type App struct {
	cfg *config.Config

	postgresClient  *postgres.Client
	redisClient     *redis.Client
	httpServer      *http.Server
	cronScheduler   *cron.Scheduler
	embeddingWorker *embedding.Worker
}

func NewApp(cfg *config.Config) (*App, error) {
	postgresClient, err := postgres.NewClient(cfg.ToPostgresConfig())
	if err != nil {
		return nil, fmt.Errorf("creating PostgreSQL client: %w", err)
	}

	jwtClient, err := jwt.NewClient(cfg.ToJWTConfig())
	if err != nil {
		return nil, fmt.Errorf("creating JWT client: %w", err)
	}

	userRepository := postgres.NewUserRepository(postgresClient)
	memoRepository := postgres.NewMemoRepository(postgresClient)
	tagRepository := postgres.NewTagRepository(postgresClient)
	collaborationRepository := postgres.NewCollaborationRepository(postgresClient)

	redisClient := redis.NewClient(cfg.ToRedisConfig())
	kvRepository := redis.NewKVRepository(redisClient)

	googleClient := google.NewClient(cfg.ToGoogleClientConfig())

	imageerClient, err := imageer.NewClient(cfg.ToImageerConfig())
	if err != nil {
		return nil, fmt.Errorf("creating Imageer client: %w", err)
	}

	var embeddingEnqueuer port.EmbeddingEnqueuer
	var embeddingRepository port.EmbeddingRepository
	var embeddingWorker *embedding.Worker
	if cfg.Embedding.Enabled {
		embeddingClient, err := embedding.NewClient(cfg.ToEmbeddingConfig())
		if err != nil {
			return nil, fmt.Errorf("creating embedding client: %w", err)
		}

		embeddingWorker = embedding.NewWorker(cfg.ToEmbeddingConfig(), embeddingClient)
		embeddingEnqueuer = embeddingWorker
		embeddingRepository = embeddingClient
	} else {
		embeddingEnqueuer = embedding.NoopEnqueuer{}
		embeddingRepository = embedding.NoopRepository{}
	}

	authService := auth.NewService(
		cfg.ToAuthServiceConfig(), postgresClient, kvRepository, userRepository, googleClient, jwtClient)
	memoService := memo.NewService(
		postgresClient, memoRepository, tagRepository, userRepository, collaborationRepository,
		embeddingEnqueuer, embeddingRepository)
	imageService := image.NewService(cfg.ToImageServiceConfig(), imageerClient)

	pingers := []port.Pinger{
		postgresClient,
		redisClient,
	}

	httpServer := http.NewServer(cfg.ToHTTPConfig(), authService, memoService, pingers, imageService)

	cronScheduler, err := cron.NewScheduler(cfg.ToCronConfig(), memoService)
	if err != nil {
		return nil, fmt.Errorf("creating cron scheduler: %w", err)
	}

	return &App{
		cfg: cfg,

		postgresClient:  postgresClient,
		redisClient:     redisClient,
		httpServer:      httpServer,
		cronScheduler:   cronScheduler,
		embeddingWorker: embeddingWorker,
	}, nil
}

func (a *App) Initialize() (err error) {
	slog.Info("component initialization start", "timeout", a.cfg.Wire.InitializeTimeout)
	start := time.Now()
	defer func() {
		slog.Info("component initialization done", "elapsed", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Wire.InitializeTimeout)
	defer cancel()

	slog.Info("migrate schemas")
	if err := a.postgresClient.MigrateSchemas(ctx); err != nil {
		return fmt.Errorf("migrating schemas: %w", err)
	}

	if a.embeddingWorker != nil {
		slog.Info("ensure qdrant collection")
		if err := a.embeddingWorker.EnsureCollection(ctx); err != nil {
			return fmt.Errorf("ensuring qdrant collection: %w", err)
		}
	}

	return nil
}

func (a *App) Run() {
	httpServerErrs := a.httpServer.Run()
	cronSchedulerErrs := a.cronScheduler.Run()

	var embeddingWorkerErrs <-chan error
	if a.embeddingWorker != nil {
		embeddingWorkerErrs = a.embeddingWorker.Run()
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-signals:
		slog.Info("received signal", "signal", s.String())
	case err := <-httpServerErrs:
		slog.Error("fatal error from http server", "error", err)
	case err := <-cronSchedulerErrs:
		slog.Error("fatal error from cron scheduler", "error", err)
	case err := <-embeddingWorkerErrs:
		slog.Error("fatal error from embedding worker", "error", err)
	}
}

func (a *App) Shutdown() {
	slog.Info("graceful shutdown start", "timeout", a.cfg.Wire.ShutdownTimeout)
	start := time.Now()
	defer func() {
		slog.Info("graceful shutdown done", "elapsed", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Wire.ShutdownTimeout)
	defer cancel()

	slog.Info("shutdown httpServer")
	if err := a.httpServer.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown httpServer", "error", err)
	}

	slog.Info("shutdown cronScheduler")
	if err := a.cronScheduler.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown cronScheduler", "error", err)
	}

	if a.embeddingWorker != nil {
		slog.Info("shutdown embeddingWorker")
		if err := a.embeddingWorker.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown embeddingWorker", "error", err)
		}
	}

	slog.Info("shutdown redisClient")
	if err := a.redisClient.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown redisClient", "error", err)
	}

	slog.Info("shutdown postgresClient")
	if err := a.postgresClient.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown postgresClient", "error", err)
	}
}
