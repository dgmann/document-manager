package http

import (
	"fmt"
	"github.com/cskr/pubsub"
	"github.com/dgmann/document-manager/api/app"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"
)

type EventService struct {
	ps            *pubsub.PubSub
	modTimeReader app.ModTimeReader
}

func NewEventService(modTimeReader app.ModTimeReader) *EventService {
	return &EventService{ps: pubsub.New(300), modTimeReader: modTimeReader}
}

func (e *EventService) Log() {
	go func() {
		c := e.Subscribe(app.EventCreated, app.EventDeleted, app.EventUpdated)
		for event := range c {
			e := event.(app.Event)
			logrus.WithFields(logrus.Fields{
				"Type":      e.Type,
				"Timestamp": e.Timestamp,
				"data":      fmt.Sprintf("%+v\n", e.Data),
			}).Info("New Event")
		}
	}()
}

func (e *EventService) Send(t app.EventType, data interface{}) {
	payload := data
	switch data.(type) {
	case *app.Record:
		payload = SetURLForRecord(data.(*app.Record), url.URL{}, e.modTimeReader)
	case []app.Record:
		payload = SetURLForRecordList(data.([]app.Record), url.URL{}, e.modTimeReader)
	}
	event := app.Event{
		Type:      t,
		Timestamp: time.Now(),
		Data:      payload,
	}
	e.ps.Pub(event, string(t))
}

func (e *EventService) Subscribe(t ...app.EventType) chan interface{} {
	types := make([]string, len(t))
	for i, et := range t {
		types[i] = string(et)
	}
	return e.ps.Sub(types...)
}
