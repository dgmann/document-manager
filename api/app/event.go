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
	Send(event Event)
}

type Event struct {
	Type      EventType `json:"type"`
	Topic     Topic     `json:"topic"`
	Timestamp time.Time `json:"timestamp"`
	Id        string    `json:"id"`
}

type Topic string

const RecordTopic = "records"

func NewEvent(topic Topic, eventType EventType, id string) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Id:        id,
		Topic:     topic,
	}
}
