package client

type WClient[T any] struct {
	*Client[T]

	write chan T
}

func NewWClient[T any]() *WClient[T] {
	return newWClient[T](nil)
}

func newWClient[T any](client *Client[T]) *WClient[T] {
	if client == nil {
		client = NewClient[T]()
	}

	wClient := &WClient[T]{
		Client: client,
		write:  make(chan T),
	}

	if err := wClient.OnClose(func() {
		close(wClient.write)
	}); err != nil {
		panic(err)
	}

	return wClient
}

// Write processes outgoing messages to SendChan.
func (c *WClient[T]) Write(msg T) error {
	if c.closed.Load() {
		return ErrClientClosed
	}
	select {
	case c.write <- msg:
	default:
		return ErrClientWriteChanFull
	}
	return nil
}

// WriteChan returns a channel to send messages to the client.
func (c *WClient[T]) WriteChan() <-chan T {
	if c.closed.Load() {
		return nil
	}
	return c.write
}
