package model

type ChatAction int

const (
	ChatActionAddMessage ChatAction = iota
	ChatActionEditMessage
	ChatActionDeleteMessage
)
