package model

import "time"

type ChatMessage struct {
	Action    ChatAction `json:"action"`
	MessageID uint64     `json:"message_id"`
	UserID    uint64     `json:"user_id"`
	Username  string     `json:"username"`
	Message   string     `json:"message"`
	Timestamp time.Time  `json:"timestamp"`
}
