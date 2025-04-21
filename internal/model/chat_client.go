package model

type ChatClient struct {
	*WsClient[ChatMessage]

	userID uint64
}

func NewChatClient(userID uint64) *ChatClient {
	return &ChatClient{
		WsClient: NewWsClient[ChatMessage](),
		userID:   userID,
	}
}

func (c *ChatClient) UserID() uint64 {
	return c.userID
}
