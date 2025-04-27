package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"stavki/internal/model"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	r *redis.Client
}

const (
	eventCacheTTL = 10 * time.Second
	userCacheTTL  = time.Minute
)

func NewCache(r *redis.Client) *Cache {
	return &Cache{
		r: r,
	}
}

func (c *Cache) R() *redis.Client {
	return c.r
}

func (c *Cache) SetEvent(ctx context.Context, event *model.Event) error {
	return Set(ctx, c.r, eventCacheKey(event.ID), event, eventCacheTTL)
}

func (c *Cache) Event(ctx context.Context, id uint64) (*model.Event, error) {
	event, err := Get[model.Event](ctx, c.r, eventCacheKey(id))
	return &event, err
}

func (c *Cache) SetUser(ctx context.Context, user *model.User) error {
	return Set(ctx, c.r, userCacheKey(user.ID), user, userCacheTTL)
}

func (c *Cache) User(ctx context.Context, id uint64) (*model.User, error) {
	event, err := Get[model.User](ctx, c.r, userCacheKey(id))
	return &event, err
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

func eventCacheKey(id uint64) string {
	return fmt.Sprintf("event:%d", id)
}

func userCacheKey(id uint64) string {
	return fmt.Sprintf("user:%d", id)
}
