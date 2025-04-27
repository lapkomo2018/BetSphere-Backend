package model

import (
	"fmt"
	"time"
)

type (
	Outcome struct {
		ID        uint64      `gorm:"primary_key" json:"id"`
		MarketID  uint64      `json:"market_id"`
		Type      OutcomeType `json:"type"`
		Price     float64     `json:"price"` // Current
		Liquidity float64     `json:"liquidity"`
		UpdatedAt time.Time   `json:"updated_at"`
		CreatedAt time.Time   `json:"created_at"`
	}

	OutcomeHistory struct {
		ID        uint64    `json:"id"`
		OutcomeID uint64    `json:"outcome_id"`
		Price     float64   `json:"price"`
		Liquidity float64   `json:"liquidity"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func NewOutcome(t OutcomeType, price, liquidity float64) *Outcome {
	return &Outcome{
		Type:      t,
		Price:     price,
		Liquidity: liquidity,
	}
}

func (o *Outcome) Print() {
	println("Outcome ID:", o.ID)
	println("Market ID:", o.MarketID)
	println("Type:", o.Type)
	fmt.Printf("Price: %.2f\n", o.Price)
	fmt.Printf("Liquidity: %.2f\n", o.Liquidity)
}
