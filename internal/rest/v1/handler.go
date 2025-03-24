package v1

import "github.com/gin-gonic/gin"

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Init(group *gin.RouterGroup) {

}
