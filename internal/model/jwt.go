package model

import (
	"time"
)

type JWT struct {
	RefreshToken string `gorm:"primaryKey"`
	UserID       uint64
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type TokenPair struct {
	AccessToken  JWTToken `json:"access_token"`
	RefreshToken JWTToken `json:"refresh_token"`
}

type JWTToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
