package service

import (
	"context"
	"fmt"
	"time"

	"stavki/internal/cache"
	"stavki/internal/database"
	"stavki/internal/log"
	"stavki/internal/model"
	"stavki/internal/model/client"
	"stavki/internal/model/event"
)

type EventService struct {
	txProvider  database.TransactionProvider
	eventDB     *database.EventRepository
	r           *cache.Cache
	userService *UserService

	hubs map[uint64]*event.OddsHub
}

// NewEvent creates a new EventService service instance.
func NewEvent(txProvider database.TransactionProvider, eventDB *database.EventRepository, r *cache.Cache, userService UserService) *EventService {
	return &EventService{
		txProvider:  txProvider,
		eventDB:     eventDB,
		r:           r,
		userService: &userService,
		hubs:        make(map[uint64]*event.OddsHub),
	}
}

func (e *EventService) Create(ctx context.Context, event *model.Event) (*model.Event, error) {
	return event, e.txProvider.Transact(func(a database.Adapters) error {
		if err := a.EventRepository.Create(ctx, event); err != nil {
			return err
		}

		for _, market := range event.Markets {
			if _, err := a.MarketRepository.CreateHistory(ctx, market); err != nil {
				return err
			}
			for _, outcome := range market.Outcomes {
				if _, err := a.OutcomeRepository.CreateHistory(ctx, outcome); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (e *EventService) Get(ctx context.Context, id uint64) (*model.Event, error) {
	return cache.Run[*model.Event](
		ctx,
		e.r.R(),
		cache.EventCacheKey(id),
		cache.EventCacheTTL,
		func(event *model.Event) (*model.Event, error) {
			if event != nil {
				return event, nil
			}

			return e.eventDB.Get(ctx, id)
		})
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
		log.WithFields(log.Fields{
			"error":  err,
			"offset": offset,
			"limit":  limit,
		}).Error("Error caching events list")
	}
	return events, nil
}

func (e *EventService) HandleBet(ctx context.Context, bet *model.Bet) error {
	var event *model.Event
	if err := e.txProvider.Transact(func(a database.Adapters) error {
		if _, err := a.UserRepository.AdjustBalance(ctx, bet.UserID, -bet.TotalAmount()); err != nil {
			return err
		}

		if err := a.BetRepository.Create(ctx, bet); err != nil {
			return err
		}

		if _, err := a.OutcomeRepository.AdjustLiquidity(ctx, bet.OutcomeID, bet.Amount); err != nil {
			return err
		}

		var err error
		event, err = a.EventRepository.Get(ctx, bet.EventID)
		if err != nil {
			return err
		}

		event.UpdateChances()
		for _, market := range event.Markets {
			if _, err := a.MarketRepository.UpdateChance(ctx, market); err != nil {
				return err
			}
			for _, outcome := range market.Outcomes {
				if _, err := a.OutcomeRepository.UpdatePrice(ctx, outcome); err != nil {
					return err
				}
			}
		}

		if err := e.r.SetEvent(ctx, event); err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"id":    event.ID,
			}).Error("Error caching event")
		}
		return nil
	}); err != nil {
		return err
	}

	e.broadcastOddsUpdate(event)
	return nil
}

func (e *EventService) ConnectOddsClient(ctx context.Context, eventID uint64, client *client.WClient[event.OddsMessage]) error {
	if _, exists := e.hubs[eventID]; !exists {
		e.hubs[eventID] = event.NewOddsHub(eventID)
	}

	hub := e.hubs[eventID]
	hub.AddClient(client)

	return client.OnClose(func() {
		hub.RemoveClient(client)
	})
}

func (e *EventService) broadcastOddsUpdate(event *model.Event) {
	hub, exists := e.hubs[event.ID]
	if !exists {
		log.WithField("event_id", event.ID).Info("odds hub not found, skipping broadcast")
		return
	}

	hub.SendMessage(event.Markets)
}
