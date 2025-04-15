package model

import (
	"fmt"
	"time"
)

type (
	Market struct {
		ID          uint64      `json:"id"`
		EventID     uint64      `json:"event_id"`
		Title       string      `json:"title"`
		Chances     float64     `json:"chances"`
		Description string      `json:"description"`
		StartTime   time.Time   `json:"start_time"`
		EndTime     time.Time   `json:"end_time"`
		Status      string      `json:"status"`   // open | closed | resolved
		Outcomes    [2]*Outcome `json:"outcomes"` // preload
	}

	MarketChancesHistory struct {
		ID        uint64    `json:"id"`
		MarketID  uint64    `json:"market_id"`
		Chances   float64   `json:"chances"`
		Timestamp time.Time `json:"timestamp"`
	}
)

func (m *Market) HandleBet(bet Bet) {
	var outcome *Outcome
	for _, o := range m.Outcomes {
		if o.ID == bet.OutcomeID {
			outcome = o
			break
		}
	}

	if outcome != nil {
		outcome.Liquidity += bet.Amount
	}
}

func (m *Market) UpdateChances() {
	total := m.Liquidity()
	if total == 0 {
		m.Chances = 0.5
	} else {
		m.Chances = m.Outcomes[0].Liquidity / total
	}
	m.UpdatePrices()
}

func (m *Market) UpdatePrices() {
	m.Outcomes[0].Price = m.Chances
	m.Outcomes[1].Price = 1 - m.Chances
}

func (m *Market) Liquidity() float64 {
	return m.Outcomes[0].Liquidity + m.Outcomes[1].Liquidity
}

func (m *Market) Print() {
	println("Market ID:", m.ID)
	println("Event ID:", m.EventID)
	println("Title:", m.Title)
	println("Description:", m.Description)
	fmt.Printf("Chances: %.2f\n", m.Chances)
	println("Start Time:", m.StartTime.String())
	println("End Time:", m.EndTime.String())
	println("Status:", m.Status)
	for _, outcome := range m.Outcomes {
		fmt.Println("---")
		outcome.Print()
	}
}
