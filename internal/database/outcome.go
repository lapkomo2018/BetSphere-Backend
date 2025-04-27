package database

import (
	"context"
	"errors"

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

func (o *OutcomeRepository) UpdateLiquidity(ctx context.Context, id uint64, amount float64) (*model.Outcome, error) {
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
