package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	r *redis.Client
}

func NewCache(r *redis.Client) *Cache {
	return &Cache{
		r: r,
	}
}

func (c *Cache) R() *redis.Client {
	return c.r
}

func Set(ctx context.Context, r *redis.Client, key string, value any, ttl time.Duration) error {
	if r == nil {
		return fmt.Errorf("redis client is nil")
	}

	jsonStr, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Set(ctx, key, jsonStr, ttl).Err()
}

func Get[T any](ctx context.Context, r *redis.Client, key string) (T, error) {
	var value T
	if r == nil {
		return value, fmt.Errorf("redis client is nil")
	}

	jsonStr, err := r.Get(ctx, key).Result()
	if err != nil {
		return value, err
	}

	if err := json.Unmarshal([]byte(jsonStr), &value); err != nil {
		return value, err
	}

	return value, nil
}
