package http

import (
	"github.com/cskr/pubsub"
	"github.com/dgmann/document-manager/api/app"
)

type EventService struct {
	ps            *pubsub.PubSub
	modTimeReader app.ModTimeReader
}

func NewEventService(modTimeReader app.ModTimeReader) *EventService {
	return &EventService{ps: pubsub.New(300), modTimeReader: modTimeReader}
}

func (e *EventService) Send(event app.Event) {
	e.ps.Pub(event, string(event.Type))
}

func (e *EventService) Subscribe(t ...app.EventType) chan interface{} {
	types := make([]string, len(t))
	for i, et := range t {
		types[i] = string(et)
	}
	return e.ps.Sub(types...)
}
