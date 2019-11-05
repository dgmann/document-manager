package event

import (
	"github.com/cskr/pubsub"
	"github.com/dgmann/document-manager/api/storage"
)

type WebsocketEventService struct {
	ps            *pubsub.PubSub
	modTimeReader storage.ModTimeReader
}

func NewWebsocketEventService(modTimeReader storage.ModTimeReader) *WebsocketEventService {
	return &WebsocketEventService{ps: pubsub.New(300), modTimeReader: modTimeReader}
}

func (e *WebsocketEventService) Send(event Event) {
	e.ps.Pub(event, string(event.Type))
}

func (e *WebsocketEventService) Subscribe(t ...Type) chan interface{} {
	types := make([]string, len(t))
	for i, et := range t {
		types[i] = string(et)
	}
	return e.ps.Sub(types...)
}
