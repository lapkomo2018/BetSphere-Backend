package cache

import (
	"context"
	"fmt"
	"time"

	"stavki/internal/model"
)

const UserCacheTTL = 10 * time.Second

func (c *Cache) SetUser(ctx context.Context, user *model.User) error {
	return Set(ctx, c.r, UserCacheKey(user.ID), user, UserCacheTTL)
}

func (c *Cache) User(ctx context.Context, id uint64) (*model.User, error) {
	event, err := Get[model.User](ctx, c.r, UserCacheKey(id))
	return &event, err
}

func UserCacheKey(id uint64) string {
	return fmt.Sprintf("user:%d", id)
}
