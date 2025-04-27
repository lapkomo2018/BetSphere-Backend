package service

import (
	"context"

	"stavki/internal/database"
	"stavki/internal/model"
)

type BetService struct {
	betDB *database.BetRepository
}

func NewBetService(betDB *database.BetRepository) *BetService {
	return &BetService{
		betDB: betDB,
	}
}

func (b *BetService) ListByUserID(ctx context.Context, userID uint64, offset, limit int) ([]*model.Bet, error) {
	return b.betDB.ListByUserID(ctx, userID, offset, limit)
}
