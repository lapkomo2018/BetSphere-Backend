package client

type RWClient[T any] struct {
	*Client[T]
	*RClient[T]
	*WClient[T]
}

func NewRWClient[T any]() *RWClient[T] {
	client := NewClient[T]()

	return &RWClient[T]{
		Client:  client,
		RClient: newRClient[T](client),
		WClient: newWClient[T](client),
	}
}
