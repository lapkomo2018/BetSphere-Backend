package database

import (
	"context"

	"stavki/internal/model"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (m *MessageRepository) Create(ctx context.Context, message *model.Message) (*model.Message, error) {
	return message, m.db.WithContext(ctx).Create(message).Error
}

func (m *MessageRepository) Get(ctx context.Context, id uint64) (*model.Message, error) {
	var message model.Message
	return &message, m.db.WithContext(ctx).Where("id = ?", id).First(&message).Error
}

func (m *MessageRepository) Save(ctx context.Context, message *model.Message) (*model.Message, error) {
	return message, m.db.WithContext(ctx).Save(message).Error
}

func (m *MessageRepository) Delete(ctx context.Context, message *model.Message) error {
	return m.db.WithContext(ctx).Delete(message).Error
}
