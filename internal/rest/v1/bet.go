package v1

import "github.com/gin-gonic/gin"

func (h *Handler) initBets(group *gin.RouterGroup) {
	group.GET("/history", h.authMiddleware, h.getBetHistory)
}

func (h *Handler) getBetHistory(c *gin.Context) {
	offset, limit := parseOffsetLimit(c)

	bets, err := h.betService.ListByUserID(c, c.GetUint64(userIDKey), offset, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get bets"})
		return
	}

	c.JSON(200, bets)
}
