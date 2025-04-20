package model

import "time"

type Message struct {
	ID        uint64 `gorm:"primary_key" json:"id"`
	EventID   uint64
	UserID    uint64
	Message   string
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
