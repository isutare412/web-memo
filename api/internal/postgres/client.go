package postgres

import (
	"context"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"

	"github.com/isutare412/web-memo/api/internal/core/ent"
)

type Client struct {
	entClient *ent.Client
}

func NewClient(cfg Config) (*Client, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
	client, err := ent.Open("postgres", dsn, buildEntOptions(cfg.QueryLog)...)
	if err != nil {
		return nil, fmt.Errorf("opening PostgreSQL connection: %w", err)
	}

	return &Client{entClient: client}, nil
}

func (c *Client) Shutdown(ctx context.Context) error {
	errs := make(chan error)
	success := make(chan struct{})
	go func() {
		if err := c.entClient.Close(); err != nil {
			errs <- fmt.Errorf("closing ent client: %w", err)
			return
		}

		close(success)
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-success:
		return nil
	}
}

func (c *Client) MigrateSchemas(ctx context.Context) error {
	if err := c.entClient.Schema.Create(ctx); err != nil {
		return fmt.Errorf("creating schema: %w", err)
	}

	if _, err := c.entClient.ExecContext(
		ctx,
		`CREATE INDEX IF NOT EXISTS memo_tags_tag_id ON memo_tags (tag_id)`,
	); err != nil {
		return fmt.Errorf("creating index: %w", err)
	}
	return nil
}

func (c *Client) BeginTx(ctx context.Context) (ctxWithTx context.Context, commit, rollback func() error) {
	tx, err := c.entClient.Tx(ctx)
	if err != nil {
		panic(err)
	}

	ctxWithTx = context.WithValue(ctx, contextTransactionKey{}, tx.Client())

	commit = func() error {
		return tx.Commit()
	}
	rollback = func() error {
		return tx.Rollback()
	}

	return ctxWithTx, commit, rollback
}

func (c *Client) WithTx(ctx context.Context, fn func(ctxWithTx context.Context) error) error {
	ctxWithTx, commit, rollback := c.BeginTx(ctx)

	defer func() {
		if v := recover(); v != nil {
			slog.Error("panicked during transaction", "recover", v)

			if err := rollback(); err != nil {
				slog.Error("failed to rollback transaction", "error", err)
			}

			panic(v)
		}
	}()

	if ferr := fn(ctxWithTx); ferr != nil {
		if rerr := rollback(); rerr != nil {
			ferr = fmt.Errorf("%w: rolling back transaction: %v", ferr, rerr)
		}
		return ferr
	}

	if err := commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

type contextTransactionKey struct{}

func transactionClient(ctx context.Context, client *ent.Client) *ent.Client {
	if txClient, ok := ctx.Value(contextTransactionKey{}).(*ent.Client); ok {
		client = txClient
	}
	return client
}

func buildEntOptions(queryLog bool) []ent.Option {
	if !queryLog {
		return nil
	}

	return []ent.Option{
		ent.Debug(),
		ent.Log(func(args ...any) {
			slog.Debug("query database", "ent", args[0])
		}),
	}
}
