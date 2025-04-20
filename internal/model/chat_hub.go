package model

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type ChatHub struct {
	eventID uint64
	clients map[*ChatClient]bool

	mu sync.RWMutex
}

func NewChatHub(eventID uint64) *ChatHub {
	return &ChatHub{
		eventID: eventID,
		clients: make(map[*ChatClient]bool),
	}
}

func (h *ChatHub) EventID() uint64 {
	return h.eventID
}

func (h *ChatHub) AddClient(client *ChatClient) {
	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()
}

func (h *ChatHub) RemoveClient(client *ChatClient) {
	h.mu.Lock()
	delete(h.clients, client)
	h.mu.Unlock()
}

func (h *ChatHub) SendMessage(msg ChatMessage) {
	h.mu.RLock()
	clients := make([]*ChatClient, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	for _, client := range clients {
		if err := client.SendMessage(msg); err != nil {
			logrus.WithFields(logrus.Fields{
				"userID":  client.UserID,
				"eventID": h.eventID,
				"error":   err,
			}).Error("failed to send message to client. closing client")
			client.Close()
			h.RemoveClient(client)
		}
	}
}
