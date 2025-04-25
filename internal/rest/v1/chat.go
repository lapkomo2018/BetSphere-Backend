package v1

import (
	"net/http"
	"strconv"

	"stavki/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	userIDKey = "userID"
)

var chatUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) initChat(group *gin.RouterGroup) {
	group.GET("/:eventID/chat", h.authMiddleware, h.chatConn)
}

func (h *Handler) chatConn(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("eventID"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Event ID"})
		return
	}

	userID := c.GetUint64(userIDKey)
	if userID == 0 {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	client := model.NewChatClient(userID)
	if err := h.chatService.HandleChatConnection(c, client, eventID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to handle chat connection"})
		return
	}

	conn, err := chatUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	if err := client.OnClose(func() {
		conn.Close()
	}); err != nil {
		c.JSON(500, gin.H{"error": "Failed to set close handler"})
		conn.Close()
		return
	}

	// Write a message to the client
	go func() {
		defer client.Close()

		for {
			select {
			case <-client.Done():
				return
			case msg := <-client.SendChan():
				if err := conn.WriteJSON(msg); err != nil {
					return
				}
			}
		}
	}()

	// Read messages from the client
	go func() {
		defer client.Close()

		for {
			var msg model.ChatMessage
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}

			if err := client.Handle(msg); err != nil {
				return
			}
		}
	}()
}
