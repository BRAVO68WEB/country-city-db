package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ListTTL   = 5 * time.Minute
	EntityTTL = 10 * time.Minute
)

func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	addr := os.Getenv("REDIS_URL")
	if addr == "" {
		addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("unable to connect to redis: %w", err)
	}

	return client, nil
}

func Get[T any](ctx context.Context, client *redis.Client, key string) (*T, error) {
	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var result T
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func Set(ctx context.Context, client *redis.Client, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.Set(ctx, key, data, ttl).Err()
}

func Delete(ctx context.Context, client *redis.Client, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return client.Del(ctx, keys...).Err()
}
