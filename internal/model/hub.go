package model

import (
	"sync"

	"stavki/internal/log"
	"stavki/internal/model/client"
)

type Hub[T any] struct {
	clients map[*client.WClient[T]]bool

	mu sync.RWMutex
}

func NewHub[T any]() *Hub[T] {
	return &Hub[T]{
		clients: make(map[*client.WClient[T]]bool),
	}
}

func (h *Hub[T]) AddClient(c *client.WClient[T]) {
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()
}

func (h *Hub[T]) RemoveClient(c *client.WClient[T]) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
}

func (h *Hub[T]) SendMessage(msg T) {
	h.mu.RLock()
	clients := make([]*client.WClient[T], 0, len(h.clients))
	for c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	for _, c := range clients {
		if err := c.Write(msg); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("failed to send message to client. closing client")
			c.Close()
			h.RemoveClient(c)
		}
	}
}
