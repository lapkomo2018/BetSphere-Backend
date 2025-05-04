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
	if market, err := m.r.Market(ctx, id); err == nil {
		return market, nil
	}

	market, err := m.marketDB.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := m.r.SetMarket(ctx, market); err != nil {
		return nil, err
	}

	return market, nil
}

func (m *MarketService) History(ctx context.Context, id uint64, offset, limit int) ([]*model.MarketChancesHistory, error) {
	if marketHistory, err := m.r.MarketHistory(ctx, id, offset, limit); err == nil {
		return marketHistory, nil
	}

	markets, err := m.marketDB.History(ctx, id, offset, limit)
	if err != nil {
		return nil, err
	}

	if err := m.r.SetMarketHistory(ctx, id, markets, offset, limit); err != nil {
		return nil, err
	}

	return markets, nil
}
