package app

import (
	"time"
)

type EventType string

const (
	EventCreated EventType = "CREATE"
	EventUpdated           = "UPDATE"
	EventDeleted           = "DELETE"
)

type Subscriber interface {
	Subscribe(t ...EventType) chan interface{}
}

type Sender interface {
	Send(t EventType, data interface{})
}

type Event struct {
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}
