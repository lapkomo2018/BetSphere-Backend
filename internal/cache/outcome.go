package cache

import (
	"context"
	"fmt"
	"time"

	"stavki/internal/model"
)

const OutcomeCacheTTL = 10 * time.Second

func (c *Cache) SetOutcome(ctx context.Context, outcome *model.Outcome) error {
	return Set(ctx, c.r, OutcomeCacheKey(outcome.ID), outcome, OutcomeCacheTTL)
}

func (c *Cache) Outcome(ctx context.Context, id uint64) (*model.Outcome, error) {
	outcome, err := Get[model.Outcome](ctx, c.r, OutcomeCacheKey(id))
	return &outcome, err
}

func (c *Cache) SetOutcomeHistory(ctx context.Context, id uint64, outcomeHistory []*model.OutcomeHistory, offset, limit int) error {
	return Set(ctx, c.r, OutcomeHistoryCacheKey(id, offset, limit), outcomeHistory, OutcomeCacheTTL)
}

func (c *Cache) OutcomeHistory(ctx context.Context, id uint64, offset, limit int) ([]*model.OutcomeHistory, error) {
	outcomeHistory, err := Get[[]*model.OutcomeHistory](ctx, c.r, OutcomeHistoryCacheKey(id, offset, limit))
	return outcomeHistory, err
}

func OutcomeCacheKey(id uint64) string {
	return fmt.Sprintf("outcome:%d", id)
}

func OutcomeHistoryCacheKey(id uint64, offset, limit int) string {
	return fmt.Sprintf("outcomehistory:%d:%d_%d", id, offset, limit)
}
