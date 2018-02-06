package services

import (
	"github.com/cskr/pubsub"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type EventType string

const (
	EventCreated EventType = "CREATE"
	EventUpdated           = "UPDATE"
	EventDeleted           = "DELETE"
)

type EventService struct {
	ps *pubsub.PubSub
}

type Event struct {
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
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
	e := &EventService{ps: pubsub.New(300)}
	go func() {
		c := e.Subscribe(EventCreated, EventDeleted, EventUpdated)
		for event := range c {
			e := event.(Event)
			log.WithFields(log.Fields{
				"Type":      e.Type,
				"Timestamp": e.Timestamp,
				"Data":      e.Data,
			}).Info("New Event")
		}
	}()
	return e
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
	for i, et := range t {
		types[i] = string(et)
	}
	return e.ps.Sub(types...)
}
