package event

import "stavki/internal/model"

type (
	OddsHub struct {
		*model.Hub[OddsMessage]

		eventID uint64
	}

	OddsMessage struct {
		Markets []*marketOdds
	}

	marketOdds struct {
		ID     uint64  `json:"id"`
		Chance float64 `json:"chance"`
	}
)

func NewOddsHub(id uint64) *OddsHub {
	return &OddsHub{
		Hub:     model.NewHub[OddsMessage](),
		eventID: id,
	}
}

func (h *OddsHub) ID() uint64 {
	return h.eventID
}

func (h *OddsHub) SendMessage(markets []*model.Market) {
	odds := make([]*marketOdds, 0, len(markets))
	for _, m := range markets {
		odds = append(odds, &marketOdds{ID: m.ID, Chance: m.Chance})
	}

	h.Hub.SendMessage(OddsMessage{Markets: odds})
}
