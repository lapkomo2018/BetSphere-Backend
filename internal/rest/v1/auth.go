package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

func (h *Handler) initAuth(group *gin.RouterGroup) {
	group.POST("/login", h.login)
	group.POST("/register", h.register)
	group.POST("/refresh", h.refresh)
	group.POST("/logout", h.logout)
}

func (h *Handler) login(c *gin.Context) {
	var body struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	_, pair, err := h.userService.Login(c.Request.Context(), body.Login, body.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(200, tokenResponse{
		AccessToken:  pair.AccessToken.Token,
		RefreshToken: pair.RefreshToken.Token,
	})
	c.SetCookie("Authorization", pair.AccessToken.Token, int(pair.AccessToken.ExpiresAt.Sub(time.Now()).Seconds()), "/", "", false, true)
}

func (h *Handler) register(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
		Password string `json:"password" binding:"required,min=8"`
		Email    string `json:"email" binding:"required,email"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	user, pair, err := h.userService.Register(c.Request.Context(), body.Username, body.Email, body.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to register"})
		return
	}

	c.JSON(200, struct {
		UserID uint64 `json:"user_id"`
		tokenResponse
	}{
		UserID: user.ID,
		tokenResponse: tokenResponse{
			AccessToken:  pair.AccessToken.Token,
			RefreshToken: pair.RefreshToken.Token,
		},
	})
	c.SetCookie("Authorization", pair.AccessToken.Token, int(pair.AccessToken.ExpiresAt.Sub(time.Now()).Seconds()), "/", "", false, true)
}

func (h *Handler) refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	pair, err := h.authService.RefreshJWT(c.Request.Context(), body.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(200, tokenResponse{
		AccessToken:  pair.AccessToken.Token,
		RefreshToken: pair.RefreshToken.Token,
	})
	c.SetCookie("Authorization", pair.AccessToken.Token, int(pair.AccessToken.ExpiresAt.Sub(time.Now()).Seconds()), "/", "", false, true)
}

func (h *Handler) logout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	if err := h.authService.LogoutJWT(c.Request.Context(), body.RefreshToken); err != nil {
		c.JSON(500, gin.H{"error": "failed to logout"})
		return
	}

	c.Status(http.StatusNoContent)
}
