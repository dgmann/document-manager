package services

import (
	"github.com/cskr/pubsub"
	"sync"
	"time"
)

type EventType string

const (
	EventCreated EventType = "CREATE"
	EventUpdated           = "UPDATE"
	EventDeleted           = "Delete"
)

type EventService struct {
	ps *pubsub.PubSub
}

type Event struct {
	Type      EventType
	Timestamp time.Time
	Data      interface{}
}

var instance *EventService
var once sync.Once

func GetEventService() *EventService {
	once.Do(func() {
		instance = newEventService()
	})
	return instance
}

func newEventService() *EventService {
	return &EventService{ps: pubsub.New(300)}
}

func (e *EventService) Send(t EventType, data interface{}) {
	event := Event{
		Type:      t,
		Timestamp: time.Now(),
		Data:      data,
	}
	e.ps.Pub(event, string(t))
}

func (e *EventService) Subscribe(t ...EventType) chan interface{} {
	types := make([]string, len(t))
	for et := range t {
		types = append(types, string(et))
	}
	return e.ps.Sub(types...)
}
