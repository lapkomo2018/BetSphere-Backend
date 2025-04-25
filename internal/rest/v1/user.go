package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initUsers(group *gin.RouterGroup) {
	group.GET("/me", h.authMiddleware, h.getMe)
	group.GET("/:id", h.getUser)
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

func (h *Handler) getUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid User ID"})
		return
	}

	user, err := h.userService.Get(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(200, struct {
		ID       uint64 `json:"id"`
		Username string `json:"username"`
	}{
		ID:       user.ID,
		Username: user.Username,
	})
}
