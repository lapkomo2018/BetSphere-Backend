package v1

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) initMarkets(group *gin.RouterGroup) {
	id := group.Group("/:id", h.marketIDMiddleware)
	{
		id.GET("", h.getMarket)
		id.GET("/history", h.getMarketHistory)
	}
}

func (h *Handler) getMarket(c *gin.Context) {
	market, err := h.marketService.Get(c, c.GetUint64(marketIDKey))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get market"})
		return
	}

	c.JSON(200, market)
}

func (h *Handler) getMarketHistory(c *gin.Context) {
	offset, limit := parseOffsetLimit(c)
	history, err := h.marketService.History(c, c.GetUint64(marketIDKey), offset, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get market history"})
		return
	}

	c.JSON(200, history)
}
