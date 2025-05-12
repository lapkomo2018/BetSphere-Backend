package service

import (
	"context"

	"stavki/internal/cache"
	"stavki/internal/database"
	"stavki/internal/model"
)

type OutcomeService struct {
	txProvider database.TransactionProvider
	outcomeDB  *database.OutcomeRepository
	r          *cache.Cache
}

func NewOutcome(txProvider database.TransactionProvider, outcomeDB *database.OutcomeRepository, r *cache.Cache) *OutcomeService {
	return &OutcomeService{
		txProvider: txProvider,
		outcomeDB:  outcomeDB,
		r:          r,
	}
}

func (o *OutcomeService) Get(ctx context.Context, id uint64) (*model.Outcome, error) {
	return cache.Run[*model.Outcome](
		ctx,
		o.r.R(),
		cache.OutcomeCacheKey(id),
		cache.OutcomeCacheTTL,
		func(outcome *model.Outcome) (*model.Outcome, error) {
			if outcome != nil {
				return outcome, nil
			}

			return o.outcomeDB.Get(ctx, id)
		})
}

func (o *OutcomeService) History(ctx context.Context, id uint64, offset, limit int) ([]*model.OutcomeHistory, error) {
	return cache.Run[[]*model.OutcomeHistory](
		ctx,
		o.r.R(),
		cache.OutcomeHistoryCacheKey(id, offset, limit),
		cache.OutcomeCacheTTL,
		func(outcomeHistory []*model.OutcomeHistory) ([]*model.OutcomeHistory, error) {
			if outcomeHistory != nil {
				return outcomeHistory, nil
			}

			return o.outcomeDB.History(ctx, id, offset, limit)
		})
}
