package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	innerClient *redis.Client
}

func NewClient(cfg Config) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
	})

	return &Client{
		innerClient: client,
	}
}

func (c *Client) Shutdown(ctx context.Context) error {
	errs := make(chan error)
	success := make(chan struct{})
	go func() {
		if err := c.innerClient.Close(); err != nil {
			errs <- fmt.Errorf("closing redis client: %w", err)
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
