package v1

import (
	"stavki/internal/service"

	"github.com/gin-gonic/gin"
)

type (
	Handler struct {
		userService *service.User
		authService *service.Auth
	}
	Config struct {
		UserService *service.User
		AuthService *service.Auth
	}
)

func New(cfg Config) *Handler {
	return &Handler{
		userService: cfg.UserService,
		authService: cfg.AuthService,
	}
}

func (h *Handler) Init(group *gin.RouterGroup) {
	h.initAuth(group.Group("/auth"))
}
