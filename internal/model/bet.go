package model

import "time"

type Bet struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"user_id"`
	EventID    uint64    `json:"event_id"`
	MarketID   uint64    `json:"market_id"`
	OutcomeID  uint64    `json:"outcome_id"`
	Amount     float64   `json:"amount"`      // $
	TokenPrice float64   `json:"token_price"` // $
	Tokens     float64   `json:"tokens"`      // amount / price
	Resolved   bool      `json:"resolved"`
	Timestamp  time.Time `json:"timestamp"`
}
