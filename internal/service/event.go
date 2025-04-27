package service

import (
	"context"
	"fmt"
	"time"

	"stavki/internal/cache"
	"stavki/internal/database"
	"stavki/internal/model"

	"github.com/sirupsen/logrus"
)

type EventService struct {
	txProvider  database.TransactionProvider
	eventDB     *database.EventRepository
	r           *cache.Cache
	userService *UserService
}

// NewEvent creates a new EventService service instance.
func NewEvent(txProvider database.TransactionProvider, eventDB *database.EventRepository, r *cache.Cache, userService UserService) *EventService {
	return &EventService{
		txProvider:  txProvider,
		eventDB:     eventDB,
		r:           r,
		userService: &userService,
	}
}

func (e *EventService) Create(ctx context.Context, event *model.Event) (*model.Event, error) {
	return event, e.txProvider.Transact(func(a database.Adapters) error {
		if err := a.EventRepository.Create(ctx, event); err != nil {
			return err
		}

		return nil
	})
}

func (e *EventService) Get(ctx context.Context, id uint64) (*model.Event, error) {
	if event, err := e.r.Event(ctx, id); err == nil {
		return event, nil
	}

	event, err := e.eventDB.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := e.r.SetEvent(ctx, event); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"id":    event.ID,
		}).Error("Error caching event")
	}
	return event, nil
}

func (e *EventService) List(ctx context.Context, offset, limit int) ([]*model.Event, error) {
	cacheKey := fmt.Sprintf("events:%d:%d", offset, limit)
	getEventsListFromCache := func(ctx context.Context, offset, limit int) ([]*model.Event, error) {
		eventsIDs, err := cache.Get[[]uint64](ctx, e.r.R(), cacheKey)
		if err != nil {
			return nil, err
		}

		events := make([]*model.Event, len(eventsIDs))
		for i, id := range eventsIDs {
			event, err := e.Get(ctx, id)
			if err != nil {
				return nil, err
			}
			events[i] = event
		}

		return events, nil
	}
	if events, err := getEventsListFromCache(ctx, offset, limit); err == nil {
		return events, nil
	}

	events, err := e.eventDB.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	eventsIDs := make([]uint64, len(events))
	for i, event := range events {
		eventsIDs[i] = event.ID
	}
	if err := cache.Set(ctx, e.r.R(), cacheKey, eventsIDs, time.Minute); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"offset": offset,
			"limit":  limit,
		}).Error("Error caching events list")
	}
	return events, nil
}

func (e *EventService) HandleBet(ctx context.Context, bet *model.Bet) error {
	return e.txProvider.Transact(func(a database.Adapters) error {
		if _, err := a.UserRepository.AdjustBalance(ctx, bet.UserID, -bet.Amount); err != nil {
			return err
		}

		if err := a.BetRepository.Create(ctx, bet); err != nil {
			return err
		}

		if _, err := a.OutcomeRepository.AdjustLiquidity(ctx, bet.OutcomeID, bet.Amount); err != nil {
			return err
		}

		event, err := a.EventRepository.Get(ctx, bet.EventID)
		if err != nil {
			return err
		}

		event.UpdateChances()
		for _, market := range event.Markets {
			if err := a.MarketRepository.Save(ctx, market); err != nil {
				return err
			}
			for _, outcome := range market.Outcomes {
				if err := a.OutcomeRepository.Save(ctx, outcome); err != nil {
					return err
				}
			}
		}

		if err := e.r.SetEvent(ctx, event); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"id":    event.ID,
			}).Error("Error caching event")
		}
		return nil
	})
}
