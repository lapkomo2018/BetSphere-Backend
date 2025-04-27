package database

import (
	"context"

	"stavki/internal/model"

	"gorm.io/gorm"
)

type MarketRepository struct {
	db *gorm.DB
}

func NewMarketRepository(db *gorm.DB) *MarketRepository {
	return &MarketRepository{
		db: db,
	}
}

func (m *MarketRepository) Create(ctx context.Context, market *model.Market) error {
	return m.db.WithContext(ctx).Create(market).Error
}

func (m *MarketRepository) Get(ctx context.Context, id uint64) (*model.Market, error) {
	var market model.Market
	return &market, m.preload(m.db.WithContext(ctx).Where("id = ?", id)).First(&market).Error
}

func (m *MarketRepository) Save(ctx context.Context, market *model.Market) error {
	return m.db.WithContext(ctx).Save(market).Error
}

func (m *MarketRepository) preload(db *gorm.DB) *gorm.DB {
	return db.Preload("Outcomes")
}
