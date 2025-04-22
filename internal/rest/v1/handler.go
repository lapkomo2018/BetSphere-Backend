package v1

import (
	"stavki/internal/service"

	"github.com/gin-gonic/gin"
)

type (
	Handler struct {
		userService *service.UserService
		authService *service.Auth
		chatService *service.ChatService
	}
	Config struct {
		UserService *service.UserService
		AuthService *service.Auth
		ChatService *service.ChatService
	}
)

func New(cfg Config) *Handler {
	return &Handler{
		userService: cfg.UserService,
		authService: cfg.AuthService,
		chatService: cfg.ChatService,
	}
}

func (h *Handler) Init(group *gin.RouterGroup) {
	h.initAuth(group.Group("/auth"))
	h.initChat(group.Group(""))
}
