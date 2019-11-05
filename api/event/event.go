package event

import (
	"time"
)

type Type string

const (
	Created Type = "CREATE"
	Updated      = "UPDATE"
	Deleted      = "DELETE"
)

type Subscriber interface {
	Subscribe(t ...Type) chan interface{}
}

type Sender interface {
	Send(event Event)
}

type Event struct {
	Type      Type      `json:"type"`
	Topic     Topic     `json:"topic"`
	Timestamp time.Time `json:"timestamp"`
	Id        string    `json:"id"`
}

type Topic string

const RecordTopic = "records"

func New(topic Topic, eventType Type, id string) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Id:        id,
		Topic:     topic,
	}
}
