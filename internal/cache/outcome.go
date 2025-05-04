package cache

import (
	"context"
	"fmt"
	"time"

	"stavki/internal/model"
)

const outcomeCacheTTL = 10 * time.Second

func (c *Cache) SetOutcome(ctx context.Context, outcome *model.Outcome) error {
	return Set(ctx, c.r, outcomeCacheKey(outcome.ID), outcome, outcomeCacheTTL)
}

func (c *Cache) Outcome(ctx context.Context, id uint64) (*model.Outcome, error) {
	outcome, err := Get[model.Outcome](ctx, c.r, outcomeCacheKey(id))
	return &outcome, err
}

func (c *Cache) SetOutcomeHistory(ctx context.Context, id uint64, outcomeHistory []*model.OutcomeHistory, offset, limit int) error {
	return Set(ctx, c.r, outcomeHistoryCacheKey(id, offset, limit), outcomeHistory, outcomeCacheTTL)
}

func (c *Cache) OutcomeHistory(ctx context.Context, id uint64, offset, limit int) ([]*model.OutcomeHistory, error) {
	outcomeHistory, err := Get[[]*model.OutcomeHistory](ctx, c.r, outcomeHistoryCacheKey(id, offset, limit))
	return outcomeHistory, err
}

func outcomeCacheKey(id uint64) string {
	return fmt.Sprintf("outcome:%d", id)
}

func outcomeHistoryCacheKey(id uint64, offset, limit int) string {
	return fmt.Sprintf("outcomehistory:%d:%d_%d", id, offset, limit)
}
