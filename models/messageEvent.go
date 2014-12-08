package models

import ()

type MessageEventType int

const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
)

type MessageEvent struct {
	Type      MessageEventType
	ClientID  string
	Timestamp int //Unix timestamp seconds
	Message   string
}
