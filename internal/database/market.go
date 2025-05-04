package database

import (
	"context"
	"time"

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

// UpdateChance updates the chance of a market and creates a history record.
func (m *MarketRepository) UpdateChance(ctx context.Context, market *model.Market) (*model.Market, error) {
	db := m.db.WithContext(ctx)
	if err := db.Model(&market).Update("chance", market.Chance).Error; err != nil {
		return market, err
	}

	if _, err := m.CreateHistory(ctx, market); err != nil {
		return market, err
	}

	return market, nil
}

func (m *MarketRepository) CreateHistory(ctx context.Context, market *model.Market) (*model.MarketChancesHistory, error) {
	history := &model.MarketChancesHistory{
		MarketID:  market.ID,
		Chances:   market.Chance,
		CreatedAt: time.Now(),
	}
	return history, m.db.WithContext(ctx).Create(history).Error
}

func (m *MarketRepository) History(ctx context.Context, id uint64, offset, limit int) ([]*model.MarketChancesHistory, error) {
	var markets []*model.MarketChancesHistory
	return markets, m.db.WithContext(ctx).
		Where("market_id = ?", id).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&markets).Error
}

func (m *MarketRepository) preload(db *gorm.DB) *gorm.DB {
	return db.Preload("Outcomes")
}
