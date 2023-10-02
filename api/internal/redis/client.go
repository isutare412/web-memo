package redis

import (
	"context"

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
	return c.innerClient.Close()
}
