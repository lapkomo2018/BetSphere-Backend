package model

import (
	"time"
)

type (
	Event struct {
		ID          uint64    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		Status      Status    `json:"status"`
		Markets     []*Market `json:"markets"` // preload
	}

	Status int
)

const (
	StatusOpen Status = iota
	StatusClosed
	StatusResolved
)

func (s Status) String() string {
	return [...]string{"open", "closed", "resolved"}[s]
}

func (e *Event) Liquidity() float64 {
	var total float64
	for _, m := range e.Markets {
		total += m.Liquidity()
	}
	return total
}

func (e *Event) HandleBet(bet Bet) {
	var market *Market
	for _, m := range e.Markets {
		if m.ID == bet.OutcomeID {
			market = m
			break
		}
	}

	if market != nil {
		market.HandleBet(bet)
		e.UpdateChances()
	}
}

func (e *Event) UpdateChances() {
	if len(e.Markets) == 1 {
		e.Markets[0].UpdateChances()
	} else {
		e.updateChancesMulti()
	}
}

func (e *Event) updateChancesMulti() {
	totalLiquidity := e.Liquidity()
	if totalLiquidity == 0 {
		chance := 1.0 / float64(len(e.Markets))
		for _, m := range e.Markets {
			m.Chances = chance
			m.UpdatePrices()
		}
		return
	}

	for _, m := range e.Markets {
		m.Chances = m.Liquidity() / totalLiquidity
		m.UpdatePrices()
	}
}

func (e *Event) Print() {
	println("Event ID:", e.ID)
	println("Name:", e.Name)
	println("Description:", e.Description)
	println("Start Time:", e.StartTime.String())
	println("End Time:", e.EndTime.String())
	println("Status:", e.Status)
	for _, market := range e.Markets {
		println("---")
		market.Print()
	}
}
