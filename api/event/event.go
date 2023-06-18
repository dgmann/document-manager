package event

import (
	"context"
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
	Send(ctx context.Context, event Event) error
}

type Event struct {
	Type      Type        `json:"type"`
	Topic     Topic       `json:"topic"`
	Timestamp time.Time   `json:"timestamp"`
	Id        string      `json:"id"`
	Data      interface{} `json:"data"`
}

type Topic string

const RecordTopic = "records"

func New(topic Topic, eventType Type, id string, data interface{}) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Id:        id,
		Topic:     topic,
		Data:      data,
	}
}
