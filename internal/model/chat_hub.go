package model

type ChatHub struct {
	*WsHub[ChatMessage]

	eventID uint64
}

func NewChatHub(eventID uint64) *ChatHub {
	return &ChatHub{
		WsHub:   NewWsHub[ChatMessage](),
		eventID: eventID,
	}
}

func (h *ChatHub) EventID() uint64 {
	return h.eventID
}
