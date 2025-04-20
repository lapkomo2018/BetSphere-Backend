package model

import (
	"errors"
	"sync"
	"sync/atomic"
)

type ChatClient struct {
	UserID uint64
	inbox  chan ChatMessage
	outbox chan ChatMessage
	done   chan struct{}

	once   sync.Once
	closed atomic.Bool
}

var (
	ErrChatClientClosed = errors.New("client is closed")
)

func NewChatClient(userID uint64) *ChatClient {
	return &ChatClient{
		UserID: userID,
		inbox:  make(chan ChatMessage),
		outbox: make(chan ChatMessage),
		done:   make(chan struct{}),
	}
}

func (c *ChatClient) SendMessage(msg ChatMessage) error {
	if c.closed.Load() {
		return ErrChatClientClosed
	}
	select {
	case c.outbox <- msg:
	default:
		return errors.New("outbox is full")
	}
	return nil
}

func (c *ChatClient) SendMessageChannel() <-chan ChatMessage {
	if c.closed.Load() {
		return nil
	}
	return c.outbox
}

func (c *ChatClient) HandleMessage(msg ChatMessage) error {
	if c.closed.Load() {
		return ErrChatClientClosed
	}
	select {
	case c.inbox <- msg:
	default:
		return errors.New("inbox is full")
	}
	return nil
}

func (c *ChatClient) MessageChannel() <-chan ChatMessage {
	if c.closed.Load() {
		return nil
	}
	return c.inbox
}

func (c *ChatClient) Close() {
	c.once.Do(func() {
		close(c.done)
		close(c.inbox)
		close(c.outbox)
		c.closed.Store(true)
	})
}

func (c *ChatClient) Done() <-chan struct{} {
	return c.done
}
