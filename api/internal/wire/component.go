package wire

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isutare412/tasks/api/internal/config"
	"github.com/isutare412/tasks/api/internal/postgres"
)

type Components struct {
	cfg            *config.Config
	postgresClient *postgres.Client
}

func NewComponents(cfg *config.Config) (*Components, error) {
	postgresClient, err := postgres.NewClient(cfg.ToPostgresConfig())
	if err != nil {
		return nil, fmt.Errorf("creating PostgreSQL client: %w", err)
	}

	return &Components{
		cfg:            cfg,
		postgresClient: postgresClient,
	}, nil
}

func (c *Components) Initialize() error {
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

	slog.Info("shutdown postgresClient")
	if err := c.postgresClient.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown postgresClient", "error", err)
	}
}
