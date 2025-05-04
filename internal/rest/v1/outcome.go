package v1

import "github.com/gin-gonic/gin"

func (h *Handler) initOutcomes(group *gin.RouterGroup) {
	id := group.Group("/:id", h.outcomeIDMiddleware)
	{
		id.GET("", h.getOutcome)
		id.GET("/history", h.getOutcomeHistory)
	}
}

func (h *Handler) getOutcome(c *gin.Context) {
	outcome, err := h.outcomeService.Get(c, c.GetUint64(outcomeIDKey))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get outcome"})
		return
	}

	c.JSON(200, outcome)
}

func (h *Handler) getOutcomeHistory(c *gin.Context) {
	offset, limit := parseOffsetLimit(c)
	outcome, err := h.outcomeService.History(c, c.GetUint64(outcomeIDKey), offset, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get outcome history"})
		return
	}

	c.JSON(200, outcome)
}
