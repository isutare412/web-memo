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
	"github.com/isutare412/web-memo/api/internal/core/service/memo"
	"github.com/isutare412/web-memo/api/internal/cron"
	"github.com/isutare412/web-memo/api/internal/google"
	"github.com/isutare412/web-memo/api/internal/http"
	"github.com/isutare412/web-memo/api/internal/jwt"
	"github.com/isutare412/web-memo/api/internal/postgres"
	"github.com/isutare412/web-memo/api/internal/redis"
)

type Components struct {
	cfg *config.Config

	postgresClient *postgres.Client
	redisClient    *redis.Client
	httpServer     *http.Server
	cronScheduler  *cron.Scheduler
}

func NewComponents(cfg *config.Config) (*Components, error) {
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

	authService := auth.NewService(
		cfg.ToAuthServiceConfig(), postgresClient, kvRepository, userRepository, googleClient, jwtClient)
	memoService := memo.NewService(
		postgresClient, memoRepository, tagRepository, userRepository, collaborationRepository)

	pingers := []port.Pinger{
		postgresClient,
		redisClient,
	}

	httpServer := http.NewServer(cfg.ToHTTPConfig(), authService, memoService, pingers)

	cronScheduler, err := cron.NewScheduler(cfg.ToCronConfig(), memoService)
	if err != nil {
		return nil, fmt.Errorf("creating cron scheduler: %w", err)
	}

	return &Components{
		cfg: cfg,

		postgresClient: postgresClient,
		redisClient:    redisClient,
		httpServer:     httpServer,
		cronScheduler:  cronScheduler,
	}, nil
}

func (c *Components) Initialize() (err error) {
	slog.Info("component initialization start", "timeout", c.cfg.Wire.InitializeTimeout)
	start := time.Now()
	defer func() {
		slog.Info("component initialization done", "elapsed", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.Wire.InitializeTimeout)
	defer cancel()

	slog.Info("migrate schemas")
	if err := c.postgresClient.MigrateSchemas(ctx); err != nil {
		return fmt.Errorf("migrating schemas: %w", err)
	}

	return nil
}

func (c *Components) Run() {
	httpServerErrs := c.httpServer.Run()
	cronSchedulerErrs := c.cronScheduler.Run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-signals:
		slog.Info("received signal", "signal", s.String())
	case err := <-httpServerErrs:
		slog.Error("fatal error from http server", "error", err)
	case err := <-cronSchedulerErrs:
		slog.Error("fatal error from cron scheduler", "error", err)
	}
}

func (c *Components) Shutdown() {
	slog.Info("graceful shutdown start", "timeout", c.cfg.Wire.ShutdownTimeout)
	start := time.Now()
	defer func() {
		slog.Info("graceful shutdown done", "elapsed", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.Wire.ShutdownTimeout)
	defer cancel()

	slog.Info("shutdown cronScheduler")
	if err := c.cronScheduler.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown cronScheduler", "error", err)
	}

	slog.Info("shutdown httpServer")
	if err := c.httpServer.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown httpServer", "error", err)
	}

	slog.Info("shutdown redisClient")
	if err := c.redisClient.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown redisClient", "error", err)
	}

	slog.Info("shutdown postgresClient")
	if err := c.postgresClient.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown postgresClient", "error", err)
	}
}
