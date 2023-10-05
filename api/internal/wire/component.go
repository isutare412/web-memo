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
	"github.com/isutare412/web-memo/api/internal/core/service/auth"
	"github.com/isutare412/web-memo/api/internal/core/service/memo"
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
}

func NewComponents(cfg *config.Config) (*Components, error) {
	postgresClient, err := postgres.NewClient(cfg.ToPostgresConfig())
	if err != nil {
		return nil, fmt.Errorf("creating PostgreSQL client: %w", err)
	}

	userRepository := postgres.NewUserRepository(postgresClient)
	memoRepository := postgres.NewMemoRepository(postgresClient)
	tagRepository := postgres.NewTagRepository(postgresClient)

	redisClient := redis.NewClient(cfg.ToRedisConfig())
	kvRepository := redis.NewKVRepository(redisClient)

	googleClient := google.NewClient(cfg.ToGoogleClientConfig())
	jwtClient := jwt.NewClient()

	authService := auth.NewService(cfg.ToAuthServiceConfig(), kvRepository, userRepository, googleClient, jwtClient)
	_ = memo.NewService(postgresClient, memoRepository, tagRepository)

	httpServer := http.NewServer(cfg.ToHTTPConfig(), authService)

	return &Components{
		cfg: cfg,

		postgresClient: postgresClient,
		redisClient:    redisClient,
		httpServer:     httpServer,
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

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-signals:
		slog.Info("received signal", "signal", s.String())
	case err := <-httpServerErrs:
		slog.Error("fatal error from http server", "error", err)
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
