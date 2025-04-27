package database

import (
	"context"

	"stavki/internal/model"

	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (e *EventRepository) Create(ctx context.Context, event *model.Event) error {
	return e.db.WithContext(ctx).Create(event).Error
}

func (e *EventRepository) Get(ctx context.Context, id uint64) (*model.Event, error) {
	var event model.Event
	return &event, e.preload(e.db.WithContext(ctx).Where("id = ?", id)).First(&event).Error
}

func (e *EventRepository) List(ctx context.Context, offset, limit int) ([]*model.Event, error) {
	var events []*model.Event
	return events, e.preload(e.db.WithContext(ctx).Offset(offset).Limit(limit)).Find(&events).Error
}

func (e *EventRepository) Save(ctx context.Context, event *model.Event) error {
	return e.db.WithContext(ctx).Save(event).Error
}

func (e *EventRepository) preload(db *gorm.DB) *gorm.DB {
	return db.Preload("Markets").Preload("Markets.Outcomes")
}
