package model

import (
	"errors"
	"fmt"
	"time"
)

type (
	Market struct {
		ID          uint64       `gorm:"primary_key" json:"id"`
		EventID     uint64       `json:"event_id"`
		Title       string       `json:"title"`
		Chance      float64      `json:"chance"`
		Description string       `json:"description"`
		StartTime   time.Time    `json:"start_time"`
		EndTime     time.Time    `json:"end_time"`
		Status      Status       `json:"status"`
		Result      *OutcomeType `json:"result"`   // yes/no
		Outcomes    []*Outcome   `json:"outcomes"` // preload
		UpdatedAt   time.Time    `json:"updated_at"`
		CreatedAt   time.Time    `json:"created_at"`
	}

	MarketChancesHistory struct {
		ID        uint64    `json:"id"`
		MarketID  uint64    `json:"market_id"`
		Chances   float64   `json:"chances"`
		CreatedAt time.Time `json:"created_at"`
	}

	OutcomeType int
)

const (
	OutcomeYes OutcomeType = iota
	OutcomeNo
)

func (o OutcomeType) String() string {
	return [...]string{"yes", "no"}[o]
}

func (m *Market) Create() error {
	m.Outcomes = []*Outcome{
		NewOutcome(OutcomeYes, 0.5, 1),
		NewOutcome(OutcomeNo, 0.5, 1),
	}

	if m.Result != nil {
		return errors.New("result already set")
	}

	return nil
}

func (m *Market) UpdateChances() {
	total := m.Liquidity()
	if total == 0 {
		m.Chance = 0.5
	} else {
		m.Chance = m.Outcomes[0].Liquidity / total
	}
	m.UpdatePrices()
}

func (m *Market) UpdatePrices() {
	m.Outcomes[0].Price = m.Chance
	m.Outcomes[1].Price = 1 - m.Chance
}

func (m *Market) Liquidity() float64 {
	return m.Outcomes[0].Liquidity + m.Outcomes[1].Liquidity
}

func (m *Market) Print() {
	println("Market ID:", m.ID)
	println("Event ID:", m.EventID)
	println("Title:", m.Title)
	println("Description:", m.Description)
	fmt.Printf("Chance: %.2f\n", m.Chance)
	println("Start Time:", m.StartTime.String())
	println("End Time:", m.EndTime.String())
	println("Status:", m.Status)
	for _, outcome := range m.Outcomes {
		fmt.Println("---")
		outcome.Print()
	}
}
