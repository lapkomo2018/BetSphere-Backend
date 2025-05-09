package model

import (
	"sync"

	"stavki/internal/log"
)

type WsClientInterface[T any] interface {
	Send(msg T) error
	Close()
}

type WsHub[T any] struct {
	clients map[WsClientInterface[T]]bool

	mu sync.RWMutex
}

func NewWsHub[T any]() *WsHub[T] {
	return &WsHub[T]{
		clients: make(map[WsClientInterface[T]]bool),
	}
}

func (h *WsHub[T]) AddClient(c WsClientInterface[T]) {
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()
}

func (h *WsHub[T]) RemoveClient(c WsClientInterface[T]) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
}

func (h *WsHub[T]) SendMessage(msg T) {
	h.mu.RLock()
	clients := make([]WsClientInterface[T], 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	for _, client := range clients {
		if err := client.Send(msg); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("failed to send message to client. closing client")
			client.Close()
			h.RemoveClient(client)
		}
	}
}
