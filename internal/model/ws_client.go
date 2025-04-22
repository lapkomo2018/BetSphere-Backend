package model

import (
	"errors"
	"sync"
	"sync/atomic"
)

type WsClient[T any] struct {
	inbox  chan T
	outbox chan T
	done   chan struct{}

	onClose   []func()
	onCloseMu sync.Mutex

	once   sync.Once
	closed atomic.Bool
}

var (
	ErrWsClientClosed     = errors.New("client is closed")
	ErrWsClientInboxFull  = errors.New("inbox is full")
	ErrWsClientOutboxFull = errors.New("outbox is full")
)

func NewWsClient[T any]() *WsClient[T] {
	return &WsClient[T]{
		inbox:  make(chan T),
		outbox: make(chan T),
		done:   make(chan struct{}),
	}
}

// Send processes outgoing messages to SendChan.
func (c *WsClient[T]) Send(msg T) error {
	if c.closed.Load() {
		return ErrWsClientClosed
	}
	select {
	case c.outbox <- msg:
	default:
		return ErrWsClientOutboxFull
	}
	return nil
}

// SendChan returns a channel to send messages to the client.
func (c *WsClient[T]) SendChan() <-chan T {
	if c.closed.Load() {
		return nil
	}
	return c.outbox
}

// Handle processes incoming messages to HandleChan.
func (c *WsClient[T]) Handle(msg T) error {
	if c.closed.Load() {
		return ErrWsClientClosed
	}
	select {
	case c.inbox <- msg:
	default:
		return ErrWsClientInboxFull
	}
	return nil
}

// HandleChan returns a channel to receive messages from the client.
func (c *WsClient[T]) HandleChan() <-chan T {
	if c.closed.Load() {
		return nil
	}
	return c.inbox
}

// OnClose registers a callback to be called when the client is closed.
func (c *WsClient[T]) OnClose(fn func()) error {
	if c.closed.Load() {
		return ErrWsClientClosed
	}
	c.onCloseMu.Lock()
	defer c.onCloseMu.Unlock()
	c.onClose = append(c.onClose, fn)
	return nil
}

// Close closes the client and calls all registered OnClose callbacks.
func (c *WsClient[T]) Close() {
	c.once.Do(func() {
		close(c.done)
		close(c.inbox)
		close(c.outbox)
		c.closed.Store(true)

		c.onCloseMu.Lock()
		defer c.onCloseMu.Unlock()
		for _, fn := range c.onClose {
			fn()
		}
	})
}

// Done returns a channel that is closed when the client is closed.
func (c *WsClient[T]) Done() <-chan struct{} {
	return c.done
}
