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
	if outcome, err := o.r.Outcome(ctx, id); err == nil {
		return outcome, nil
	}

	outcome, err := o.outcomeDB.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := o.r.SetOutcome(ctx, outcome); err != nil {
		return nil, err
	}

	return outcome, nil
}

func (o *OutcomeService) History(ctx context.Context, id uint64, offset, limit int) ([]*model.OutcomeHistory, error) {
	if outcomeHistory, err := o.r.OutcomeHistory(ctx, id, offset, limit); err == nil {
		return outcomeHistory, nil
	}

	outcomes, err := o.outcomeDB.History(ctx, id, offset, limit)
	if err != nil {
		return nil, err
	}

	if err := o.r.SetOutcomeHistory(ctx, id, outcomes, offset, limit); err != nil {
		return nil, err
	}

	return outcomes, nil
}
