package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/isutare412/web-memo/api/internal/pkgerr"
	"github.com/isutare412/web-memo/api/internal/tracing"
)

type KVRepository struct {
	client *redis.Client
}

func NewKVRepository(client *Client) *KVRepository {
	return &KVRepository{
		client: client.innerClient,
	}
}

func (r *KVRepository) Get(ctx context.Context, key string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "redis.KVRepository.Get")
	defer span.End()

	val, err := r.client.Get(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return "", pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			Origin:    err,
			ClientMsg: fmt.Sprintf("key(%s) does not exist", key),
		}
	case err != nil:
		return "", fmt.Errorf("getting key: %w", err)
	}

	return val, nil
}

func (r *KVRepository) GetThenDelete(ctx context.Context, key string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "redis.KVRepository.GetThenDelete")
	defer span.End()

	val, err := r.client.GetDel(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return "", pkgerr.Known{
			Code:      pkgerr.CodeNotFound,
			Origin:    err,
			ClientMsg: fmt.Sprintf("key(%s) does not exist", key),
		}
	case err != nil:
		return "", fmt.Errorf("getting key: %w", err)
	}

	return val, nil
}

func (r *KVRepository) Set(ctx context.Context, key, val string, exp time.Duration) error {
	ctx, span := tracing.StartSpan(ctx, "redis.KVRepository.Set")
	defer span.End()

	_, err := r.client.Set(ctx, key, val, exp).Result()
	if err != nil {
		return fmt.Errorf("setting key-value: %w", err)
	}
	return nil
}

func (r *KVRepository) Delete(ctx context.Context, keys ...string) (delCount int64, err error) {
	ctx, span := tracing.StartSpan(ctx, "redis.KVRepository.Delete")
	defer span.End()

	count, err := r.client.Del(ctx, keys...).Result()
	if err != nil {
		return 0, fmt.Errorf("deleting keys: %w", err)
	}
	return count, nil
}
