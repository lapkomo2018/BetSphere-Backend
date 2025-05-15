package chat

import "stavki/internal/model"

type Hub struct {
	*model.Hub[Message]

	eventID uint64
}

func NewHub(eventID uint64) *Hub {
	return &Hub{
		Hub:     model.NewHub[Message](),
		eventID: eventID,
	}
}

func (h *Hub) EventID() uint64 {
	return h.eventID
}
