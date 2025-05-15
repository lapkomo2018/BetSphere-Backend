package client

import (
	"sync"
	"sync/atomic"
)

type Client[T any] struct {
	done chan struct{}

	onClose   []func()
	onCloseMu sync.Mutex

	once   sync.Once
	closed atomic.Bool
}

func NewClient[T any]() *Client[T] {
	return &Client[T]{
		done: make(chan struct{}),
	}
}

// OnClose registers a callback to be called when the client is closed.
func (c *Client[T]) OnClose(fn func()) error {
	if c.closed.Load() {
		return ErrClientClosed
	}
	c.onCloseMu.Lock()
	defer c.onCloseMu.Unlock()
	c.onClose = append(c.onClose, fn)
	return nil
}

// Close closes the client and calls all registered OnClose callbacks.
func (c *Client[T]) Close() {
	c.once.Do(func() {
		c.closed.Store(true)

		c.onCloseMu.Lock()
		handlers := append([]func(){}, c.onClose...)
		c.onCloseMu.Unlock()

		close(c.done)

		for _, fn := range handlers {
			fn()
		}
	})
}

// Done returns a channel that is closed when the client is closed.
func (c *Client[T]) Done() <-chan struct{} {
	return c.done
}
