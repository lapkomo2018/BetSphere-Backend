package database

import (
	"context"

	"stavki/internal/model"

	"gorm.io/gorm"
)

type BetRepository struct {
	db *gorm.DB
}

func NewBetRepository(db *gorm.DB) *BetRepository {
	return &BetRepository{
		db: db,
	}
}

func (b *BetRepository) Create(ctx context.Context, bet *model.Bet) error {
	return b.db.WithContext(ctx).Create(bet).Error
}

func (b *BetRepository) Get(ctx context.Context, id uint64) (*model.Bet, error) {
	var bet model.Bet
	return &bet, b.db.WithContext(ctx).Where("id = ?", id).First(&bet).Error
}

func (b *BetRepository) ListByUserID(ctx context.Context, userID uint64, offset, limit int) ([]*model.Bet, error) {
	var bets []*model.Bet
	err := b.db.WithContext(ctx).Offset(offset).Limit(limit).Where("user_id = ?", userID).Find(&bets).Error
	if err != nil {
		return nil, err
	}
	return bets, nil
}
