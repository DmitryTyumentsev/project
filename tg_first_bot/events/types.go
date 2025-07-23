package events

import "example.com/clients"

const (
	UnknownType EventType = iota
	Message
)

func NewEventProcess(client *clients.Client, offset int) *EventProcessor {
	return &EventProcessor{
		Client: client,
		offset: offset,
	}
}

type EventType int
type Event struct {
	Type     EventType
	TextPage string
	Meta
}

type Meta struct {
	ChatID   int
	UpdateID int
	//Delete    bool
	//UserID    int
	//MessageID int
	//Username  string
}

type EventProcessor struct {
	Client *clients.Client
	offset int
}
