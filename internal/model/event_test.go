package model

import (
	"fmt"
	"testing"
	"time"
)

func TestEvent_HandleBetMulti(t *testing.T) {
	// Create an event with markets
	event := Event{
		ID:          1,
		Name:        "Test Event",
		Description: "Test Description",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "open",
		Markets: []*Market{
			{
				ID:          1,
				EventID:     1,
				Title:       "Test Market 1",
				Description: "Test Description 1",
				StartTime:   time.Now(),
				EndTime:     time.Now().Add(24 * time.Hour),
				Status:      "open",
				Outcomes: [2]*Outcome{
					{ID: 1, MarketID: 1, Name: "Yes", Price: 0.5, Liquidity: 100},
					{ID: 2, MarketID: 1, Name: "No", Price: 0.5, Liquidity: 100},
				},
			},
			{
				ID:          2,
				EventID:     1,
				Title:       "Test Market 2",
				Description: "Test Description 2",
				StartTime:   time.Now(),
				EndTime:     time.Now().Add(24 * time.Hour),
				Status:      "open",
				Outcomes: [2]*Outcome{
					{ID: 3, MarketID: 2, Name: "Yes", Price: 0.5, Liquidity: 100},
					{ID: 4, MarketID: 2, Name: "No", Price: 0.5, Liquidity: 100},
				},
			},
		},
	}
	event.Print()
	fmt.Println()

	for range 5 {
		bet := Bet{
			ID:              1,
			EventID:         1,
			MarketID:        1,
			OutcomeID:       1,
			UserID:          1,
			Amount:          50,
			PriceAtPurchase: 0.5,
			Tokens:          100,
			Timestamp:       time.Now(),
		}

		event.HandleBet(bet)
		event.Print()
		fmt.Println()
	}
}

func TestEvent_HandleBet(t *testing.T) {
	// Create an event with markets
	event := Event{
		ID:          1,
		Name:        "Test Event",
		Description: "Test Description",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(24 * time.Hour),
		Status:      "open",
		Markets: []*Market{
			{
				ID:          1,
				EventID:     1,
				Title:       "Test Market 1",
				Description: "Test Description 1",
				StartTime:   time.Now(),
				EndTime:     time.Now().Add(24 * time.Hour),
				Status:      "open",
				Outcomes: [2]*Outcome{
					{ID: 1, MarketID: 1, Name: "Yes", Price: 0.5, Liquidity: 100},
					{ID: 2, MarketID: 1, Name: "No", Price: 0.5, Liquidity: 100},
				},
			},
		},
	}
	event.Print()
	fmt.Println()

	for range 5 {
		bet := Bet{
			ID:              1,
			EventID:         1,
			MarketID:        1,
			OutcomeID:       1,
			UserID:          1,
			Amount:          50,
			PriceAtPurchase: 0.5,
			Tokens:          100,
			Timestamp:       time.Now(),
		}

		event.HandleBet(bet)
		event.Print()
		fmt.Println()
	}
}
