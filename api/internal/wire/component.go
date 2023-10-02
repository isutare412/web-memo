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
	"github.com/isutare412/web-memo/api/internal/core/service/memo"
	"github.com/isutare412/web-memo/api/internal/postgres"
	"github.com/isutare412/web-memo/api/internal/redis"
)

type Components struct {
	cfg            *config.Config
	postgresClient *postgres.Client
	userRepository *postgres.UserRepository
	memoRepository *postgres.MemoRepository
	tagRepository  *postgres.TagRepository
	redisClient    *redis.Client
	kvRepository   *redis.KVRepository
	memoService    *memo.Service
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

	memoService := memo.NewService(postgresClient, memoRepository, tagRepository)

	return &Components{
		cfg:            cfg,
		postgresClient: postgresClient,
		userRepository: userRepository,
		memoRepository: memoRepository,
		tagRepository:  tagRepository,
		redisClient:    redisClient,
		kvRepository:   kvRepository,
		memoService:    memoService,
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
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	s := <-signals
	slog.Info("received signal", "signal", s.String())
}

func (c *Components) Shutdown() {
	slog.Info("graceful shutdown start", "timeout", c.cfg.Wire.ShutdownTimeout)
	start := time.Now()
	defer func() {
		slog.Info("graceful shutdown done", "elapsed", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.Wire.ShutdownTimeout)
	defer cancel()

	slog.Info("shutdown redisClient")
	if err := c.redisClient.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown redisClient", "error", err)
	}

	slog.Info("shutdown postgresClient")
	if err := c.postgresClient.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown postgresClient", "error", err)
	}
}
