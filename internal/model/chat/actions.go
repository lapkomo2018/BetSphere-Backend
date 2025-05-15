package chat

type Action int

const (
	ActionAddMessage Action = iota
	ActionEditMessage
	ActionDeleteMessage
)
