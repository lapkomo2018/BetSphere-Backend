package v1

import (
	"fmt"
	"math"
	"time"

	"stavki/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initEvents(group *gin.RouterGroup) {
	group.POST("", h.authMiddleware, h.requireAdmin, h.createEvent)
	group.GET("", h.getEvents)
	id := group.Group("/:id", h.eventIDMiddleware)
	{
		id.GET("", h.getEvent)
		id.GET("/chat", h.authMiddleware, h.chatConn)
		id.POST("/bet", h.authMiddleware, h.placeBet)
	}
}

func (h *Handler) createEvent(c *gin.Context) {
	var body struct {
		Name        string    `json:"name" binding:"required"`
		Description string    `json:"description" binding:"required"`
		StartTime   time.Time `json:"start_time" binding:"required"`
		EndTime     time.Time `json:"end_time" binding:"required"`
		Markets     []struct {
			Title       string `json:"title" binding:"required"`
			Description string `json:"description" binding:"required"`
		} `json:"markets" binding:"required"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	event := &model.Event{
		Name:        body.Name,
		Description: body.Description,
		StartTime:   body.StartTime,
		EndTime:     body.EndTime,
		Markets:     make([]*model.Market, len(body.Markets)),
	}
	for i, m := range body.Markets {
		event.Markets[i] = &model.Market{
			Title:       m.Title,
			Description: m.Description,
			StartTime:   body.StartTime,
			EndTime:     body.EndTime,
		}
	}
	if err := event.Create(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	event, err := h.eventService.Create(c.Request.Context(), event)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(201, event)
}

func (h *Handler) getEvent(c *gin.Context) {
	id := c.GetUint64(eventIDKey)
	if id == 0 {
		c.JSON(400, gin.H{"error": "invalid event ID"})
		return
	}

	event, err := h.eventService.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "event not found"})
		return
	}

	c.JSON(200, event)
}

func (h *Handler) getEvents(c *gin.Context) {
	offset, limit := parseOffsetLimit(c)

	events, err := h.eventService.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get events"})
		return
	}

	c.JSON(200, events)
}

func (h *Handler) placeBet(c *gin.Context) {
	var body struct {
		OutcomeID        uint64  `json:"outcome_id" binding:"required"`
		Amount           float64 `json:"amount" binding:"required"`
		TokenPrice       float64 `json:"token_price" binding:"required"`
		AllowedDeviation float64 `json:"allowed_deviation" binding:"required"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	if body.Amount <= 0 {
		c.JSON(400, gin.H{"error": "invalid amount"})
		return
	}

	if body.AllowedDeviation <= 0 {
		c.JSON(400, gin.H{"error": "invalid allowed deviation"})
		return
	}

	event, err := h.eventService.Get(c.Request.Context(), c.GetUint64(eventIDKey))
	if err != nil {
		c.JSON(404, gin.H{"error": "event not found"})
		return
	}

	if event.EndTime.Before(time.Now()) {
		c.JSON(400, gin.H{"error": "event has ended"})
		return
	}

	var outcome *model.Outcome
	var marketID uint64
	for _, market := range event.Markets {
		for _, o := range market.Outcomes {
			if o.ID == body.OutcomeID {
				outcome = o
				marketID = market.ID
				break
			}
		}
	}
	if outcome == nil {
		c.JSON(400, gin.H{"error": "outcome not found"})
		return
	}

	deviation := math.Abs(outcome.Price-body.TokenPrice) / outcome.Price * 100
	if deviation > body.AllowedDeviation {
		c.JSON(400, gin.H{"error": fmt.Sprintf("allowed deviation exceeded currently at %.2f%%", deviation)})
		return
	}

	bet := &model.Bet{
		UserID:     c.GetUint64(userIDKey),
		EventID:    event.ID,
		MarketID:   marketID,
		OutcomeID:  outcome.ID,
		Amount:     body.Amount,
		TokenPrice: outcome.Price,
		Tokens:     body.Amount / outcome.Price,
		Timestamp:  time.Now(),
	}
	if err := h.eventService.HandleBet(c.Request.Context(), bet); err != nil {
		c.JSON(500, gin.H{"error": "failed to place bet: " + err.Error()})
		return
	}

	c.JSON(201, bet)
}
