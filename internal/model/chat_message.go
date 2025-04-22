package model

import "time"

type ChatMessage struct {
	ID        uint64     `json:"id"`
	Action    ChatAction `json:"action"`
	UserID    uint64     `json:"user_id"`
	Username  string     `json:"username"`
	Message   string     `json:"message"`
	Timestamp time.Time  `json:"timestamp"`
}
