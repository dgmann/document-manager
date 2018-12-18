package services

import (
	"fmt"
	"github.com/cskr/pubsub"
	log "github.com/sirupsen/logrus"
	"time"
)

type EventType string

const (
	EventCreated EventType = "CREATE"
	EventUpdated           = "UPDATE"
	EventDeleted           = "DELETE"
)

type EventService struct {
	ps            *pubsub.PubSub
	modTimeReader ModTimeReader
}

type Event struct {
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func NewEventService(modTimeReader ModTimeReader) *EventService {
	return &EventService{ps: pubsub.New(300), modTimeReader: modTimeReader}
}

func (e *EventService) Log() {
	go func() {
		c := e.Subscribe(EventCreated, EventDeleted, EventUpdated)
		for event := range c {
			e := event.(Event)
			log.WithFields(log.Fields{
				"Type":      e.Type,
				"Timestamp": e.Timestamp,
				"data":      fmt.Sprintf("%+v\n", e.Data),
			}).Info("New Event")
		}
	}()
}

func (e *EventService) Send(t EventType, data interface{}) {
	payload := SetURL(data, "", e.modTimeReader)
	event := Event{
		Type:      t,
		Timestamp: time.Now(),
		Data:      payload,
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
