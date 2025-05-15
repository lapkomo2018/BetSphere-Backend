package client

type RClient[T any] struct {
	*Client[T]

	read chan T
}

func NewRClient[T any]() *RClient[T] {
	return newRClient[T](nil)
}

func newRClient[T any](client *Client[T]) *RClient[T] {
	if client == nil {
		client = NewClient[T]()
	}

	rClient := &RClient[T]{
		Client: client,
		read:   make(chan T),
	}

	if err := rClient.OnClose(func() {
		close(rClient.read)
	}); err != nil {
		panic(err)
	}

	return rClient
}

// Read processes incoming messages to ReadChan.
func (c *RClient[T]) Read(msg T) error {
	if c.closed.Load() {
		return ErrClientClosed
	}
	select {
	case c.read <- msg:
	default:
		return ErrClientReadChanFull
	}
	return nil
}

// ReadChan returns a channel to receive messages from the client.
func (c *RClient[T]) ReadChan() <-chan T {
	if c.closed.Load() {
		return nil
	}
	return c.read
}
