package cache

import (
	"context"
	"fmt"
	"time"

	"stavki/internal/model"
)

const EventCacheTTL = 10 * time.Second

func (c *Cache) SetEvent(ctx context.Context, event *model.Event) error {
	return Set(ctx, c.r, EventCacheKey(event.ID), event, EventCacheTTL)
}

func (c *Cache) Event(ctx context.Context, id uint64) (*model.Event, error) {
	event, err := Get[model.Event](ctx, c.r, EventCacheKey(id))
	return &event, err
}

func EventCacheKey(id uint64) string {
	return fmt.Sprintf("event:%d", id)
}
