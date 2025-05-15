package chat

import (
	"stavki/internal/model/client"
)

type Client struct {
	*client.RWClient[Message]

	userID uint64
}

func NewClient(userID uint64) *Client {
	return &Client{
		RWClient: client.NewRWClient[Message](),
		userID:   userID,
	}
}

func (c *Client) UserID() uint64 {
	return c.userID
}
