package model

import (
	"time"
)

type Bet struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"user_id"`
	EventID    uint64    `json:"event_id"`
	MarketID   uint64    `json:"market_id"`
	OutcomeID  uint64    `json:"outcome_id"`
	Amount     float64   `json:"amount"`      // $
	Fee        float64   `json:"fee"`         // $
	TokenPrice float64   `json:"token_price"` // $
	Tokens     float64   `json:"tokens"`      // amount / price
	Won        *bool     `json:"won"`         // nil if not resolved
	Resolved   bool      `json:"resolved"`
	Timestamp  time.Time `json:"timestamp"`
}

func NewBet(userID, eventID, marketID, outcomeID uint64, amount float64) *Bet {
	return &Bet{
		UserID:    userID,
		EventID:   eventID,
		MarketID:  marketID,
		OutcomeID: outcomeID,
		Amount:    amount,
		Timestamp: time.Now(),
	}
}

func (b *Bet) CalculateFee(p float64) *Bet {
	b.Fee = b.Amount * p / 100
	return b
}

func (b *Bet) CalculateTokens(price float64) *Bet {
	b.TokenPrice = price
	b.Tokens = b.Amount / b.TokenPrice
	return b
}

func (b *Bet) CalculateDeviation(expected float64) float64 {
	return (b.TokenPrice - expected) / expected * 100
}

func (b *Bet) TotalAmount() float64 {
	return b.Amount + b.Fee
}
