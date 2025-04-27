package model

import (
	"errors"
	"time"
)

type (
	Event struct {
		ID          uint64    `gorm:"primary_key" json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		Status      Status    `json:"status"`
		Markets     []*Market `json:"markets"` // preload
		UpdatedAt   time.Time `json:"updated_at"`
		CreatedAt   time.Time `json:"created_at"`
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

func (e *Event) Create() error {
	if len(e.Markets) == 0 {
		return errors.New("no Markets defined")
	}

	e.Status = StatusOpen
	if e.StartTime.IsZero() || e.StartTime.Before(time.Now()) {
		e.StartTime = time.Now()
	}

	if e.EndTime.IsZero() || e.EndTime.Before(e.StartTime) {
		return errors.New("end time must be greater than start time")
	}

	for _, m := range e.Markets {
		m.StartTime = e.StartTime
		m.EndTime = e.EndTime
		m.Status = StatusOpen
		if err := m.Create(); err != nil {
			return err
		}
	}

	e.UpdateChances()
	return nil
}

func (e *Event) Liquidity() float64 {
	var total float64
	for _, m := range e.Markets {
		total += m.Liquidity()
	}
	return total
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
			m.Chance = chance
			m.UpdatePrices()
		}
		return
	}

	for _, m := range e.Markets {
		m.Chance = m.Liquidity() / totalLiquidity
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
