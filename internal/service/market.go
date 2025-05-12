package service

import (
	"context"

	"stavki/internal/cache"
	"stavki/internal/database"
	"stavki/internal/model"
)

type MarketService struct {
	txProvider database.TransactionProvider
	marketDB   *database.MarketRepository
	r          *cache.Cache
}

func NewMarket(txProvider database.TransactionProvider, marketDB *database.MarketRepository, r *cache.Cache) *MarketService {
	return &MarketService{
		txProvider: txProvider,
		marketDB:   marketDB,
		r:          r,
	}
}

func (m *MarketService) Get(ctx context.Context, id uint64) (*model.Market, error) {
	return cache.Run[*model.Market](
		ctx,
		m.r.R(),
		cache.MarketCacheKey(id),
		cache.MarketCacheTTL,
		func(market *model.Market) (*model.Market, error) {
			if market != nil {
				return market, nil
			}

			return m.marketDB.Get(ctx, id)
		})
}

func (m *MarketService) History(ctx context.Context, id uint64, offset, limit int) ([]*model.MarketChancesHistory, error) {
	return cache.Run[[]*model.MarketChancesHistory](
		ctx,
		m.r.R(),
		cache.MarketHistoryCacheKey(id, offset, limit),
		cache.MarketCacheTTL,
		func(marketHistory []*model.MarketChancesHistory) ([]*model.MarketChancesHistory, error) {
			if marketHistory != nil {
				return marketHistory, nil
			}

			return m.marketDB.History(ctx, id, offset, limit)
		})
}
