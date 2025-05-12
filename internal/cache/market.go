package cache

import (
	"context"
	"fmt"
	"time"

	"stavki/internal/model"
)

const MarketCacheTTL = 10 * time.Second

func (c *Cache) SetMarket(ctx context.Context, market *model.Market) error {
	return Set(ctx, c.r, MarketCacheKey(market.ID), market, MarketCacheTTL)
}

func (c *Cache) Market(ctx context.Context, id uint64) (*model.Market, error) {
	market, err := Get[model.Market](ctx, c.r, MarketCacheKey(id))
	return &market, err
}

func (c *Cache) SetMarketHistory(ctx context.Context, id uint64, marketHistory []*model.MarketChancesHistory, offset, limit int) error {
	return Set(ctx, c.r, MarketHistoryCacheKey(id, offset, limit), marketHistory, MarketCacheTTL)
}

func (c *Cache) MarketHistory(ctx context.Context, id uint64, offset, limit int) ([]*model.MarketChancesHistory, error) {
	marketHistory, err := Get[[]*model.MarketChancesHistory](ctx, c.r, MarketHistoryCacheKey(id, offset, limit))
	if err != nil {
		return nil, err
	}
	return marketHistory, nil
}

func MarketCacheKey(id uint64) string {
	return fmt.Sprintf("market:%d", id)
}

func MarketHistoryCacheKey(id uint64, offset, limit int) string {
	return fmt.Sprintf("markethistory:%d:%d_%d", id, offset, limit)
}
