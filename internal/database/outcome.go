package database

import (
	"context"
	"errors"
	"time"

	"stavki/internal/model"

	"gorm.io/gorm"
)

type OutcomeRepository struct {
	db *gorm.DB
}

func NewOutcomeRepository(db *gorm.DB) *OutcomeRepository {
	return &OutcomeRepository{
		db: db,
	}
}

func (o *OutcomeRepository) Create(ctx context.Context, outcome *model.Outcome) error {
	return o.db.WithContext(ctx).Create(outcome).Error
}

func (o *OutcomeRepository) Get(ctx context.Context, id uint64) (*model.Outcome, error) {
	var outcome model.Outcome
	return &outcome, o.db.WithContext(ctx).Where("id = ?", id).First(&outcome).Error
}

func (o *OutcomeRepository) Save(ctx context.Context, outcome *model.Outcome) error {
	return o.db.WithContext(ctx).Save(outcome).Error
}

func (o *OutcomeRepository) AdjustLiquidity(ctx context.Context, id uint64, amount float64) (*model.Outcome, error) {
	var outcome model.Outcome
	result := o.db.WithContext(ctx).
		Model(&outcome).
		Where("id = ?", id).
		Where("liquidity + ? >= 0", amount).
		Update("liquidity", gorm.Expr("liquidity + ?", amount))

	if result.Error != nil {
		return &outcome, result.Error
	}

	if result.RowsAffected == 0 {
		return &outcome, errors.New("insufficient liquidity")
	}

	return &outcome, nil
}

// UpdatePrice updates the price of an outcome and creates a history record.
func (o *OutcomeRepository) UpdatePrice(ctx context.Context, outcome *model.Outcome) (*model.Outcome, error) {
	db := o.db.WithContext(ctx)
	if err := db.Model(&outcome).Update("price", outcome.Price).Error; err != nil {
		return outcome, err
	}

	if _, err := o.CreateHistory(ctx, outcome); err != nil {
		return outcome, err
	}

	return outcome, nil
}

func (o *OutcomeRepository) CreateHistory(ctx context.Context, outcome *model.Outcome) (*model.OutcomeHistory, error) {
	history := &model.OutcomeHistory{
		OutcomeID: outcome.ID,
		Price:     outcome.Price,
		Liquidity: outcome.Liquidity,
		CreatedAt: time.Now(),
	}
	return history, o.db.WithContext(ctx).Create(history).Error
}
