package v1

import (
	"net/http"

	"stavki/internal/model/chat"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var chatUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) initChat(group *gin.RouterGroup) {
}

func (h *Handler) chatConn(c *gin.Context) {
	eventID := c.GetUint64(eventIDKey)
	if eventID == 0 {
		c.JSON(400, gin.H{"error": "invalid event ID"})
		return
	}

	userID := c.GetUint64(userIDKey)
	if userID == 0 {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	client := chat.NewClient(userID)
	if err := h.chatService.HandleChatConnection(c, client, eventID); err != nil {
		c.JSON(500, gin.H{"error": "failed to handle chat connection"})
		return
	}

	conn, err := chatUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to upgrade connection"})
		return
	}

	if err := client.OnClose(func() {
		conn.Close()
	}); err != nil {
		c.JSON(500, gin.H{"error": "failed to set close handler"})
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
			case msg := <-client.WriteChan():
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
			var msg chat.Message
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}

			if err := client.Read(msg); err != nil {
				return
			}
		}
	}()
}
