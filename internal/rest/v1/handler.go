package v1

import (
	"stavki/external/validate"
	"stavki/internal/service"

	"github.com/gin-gonic/gin"
)

type (
	Handler struct {
		userService    *service.UserService
		authService    *service.Auth
		chatService    *service.ChatService
		eventService   *service.EventService
		betService     *service.BetService
		marketService  *service.MarketService
		outcomeService *service.OutcomeService
	}
	Config struct {
		UserService    *service.UserService
		AuthService    *service.Auth
		ChatService    *service.ChatService
		EventService   *service.EventService
		BetService     *service.BetService
		MarketService  *service.MarketService
		OutcomeService *service.OutcomeService
	}
)

func New(cfg Config) (*Handler, error) {
	return &Handler{
		userService:    cfg.UserService,
		authService:    cfg.AuthService,
		chatService:    cfg.ChatService,
		eventService:   cfg.EventService,
		betService:     cfg.BetService,
		marketService:  cfg.MarketService,
		outcomeService: cfg.OutcomeService,
	}, validate.StructPointersNotNil(cfg)
}

func (h *Handler) Init(group *gin.RouterGroup) {
	h.initAuth(group.Group("/auth"))
	h.initUsers(group.Group("/users"))
	h.initEvents(group.Group("/events"))
	h.initBets(group.Group("/bets"))
	h.initChat(group.Group(""))
	h.initMarkets(group.Group("/markets"))
	h.initOutcomes(group.Group("/outcomes"))
}
