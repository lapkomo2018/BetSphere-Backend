package v1

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	userIDKey    = "userID"
	eventIDKey   = "eventID"
	marketIDKey  = "marketID"
	outcomeIDKey = "outcomeID"
)

// authMiddleware is a middleware that checks if the user is authenticated
// and sets the user ID in the context
func (h *Handler) authMiddleware(c *gin.Context) {
	var token string
	if t := c.Request.Header.Get("Authorization"); t != "" {
		token = t
	} else if t, err := c.Cookie("Authorization"); err == nil {
		token = t
	} else if t := c.Query("token"); t != "" {
		token = t
	} else {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}

	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := h.authService.AuthenticateJWT(c, strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	c.Set(userIDKey, userID)
}

func (h *Handler) requireAdmin(c *gin.Context) {
	userID := c.GetUint64(userIDKey)
	if userID == 0 {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.Get(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "failed to get user"})
		return
	}

	if !user.Admin {
		c.AbortWithStatus(403)
		return
	}
}

func (h *Handler) eventIDMiddleware(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid Event ID"})
		return
	}

	if _, err := h.eventService.Get(c.Request.Context(), id); err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "event not found"})
		return
	}

	c.Set(eventIDKey, id)
}

func (h *Handler) marketIDMiddleware(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid Market ID"})
		return
	}

	if _, err := h.marketService.Get(c.Request.Context(), id); err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "market not found"})
		return
	}

	c.Set(marketIDKey, id)
}

func (h *Handler) outcomeIDMiddleware(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid Outcome ID"})
		return
	}

	if _, err := h.outcomeService.Get(c.Request.Context(), id); err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "outcome not found"})
		return
	}

	c.Set(outcomeIDKey, id)
}

func parseOffsetLimit(c *gin.Context) (offset, limit int) {
	var body struct {
		Offset int `form:"offset"`
		Limit  int `form:"limit"`
	}
	if err := c.BindQuery(&body); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"query": c.Request.URL.RawQuery,
		}).Error("failed to parse query, using default values")
	}
	if body.Offset <= 0 {
		body.Offset = 0
	}
	if body.Limit <= 0 {
		body.Limit = 10
	}

	return body.Offset, body.Limit
}
