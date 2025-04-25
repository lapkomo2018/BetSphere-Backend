package v1

import "github.com/gin-gonic/gin"

func (h *Handler) initUsers(group *gin.RouterGroup) {
	group.GET("/me", h.authMiddleware, h.getMe)
}

func (h *Handler) getMe(c *gin.Context) {
	userID := c.GetUint64(userIDKey)
	if userID == 0 {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.Get(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(200, user)
}
