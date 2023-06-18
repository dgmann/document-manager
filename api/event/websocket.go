package event

import (
	"context"
	"github.com/cskr/pubsub"
)

type WebsocketEventService struct {
	ps *pubsub.PubSub
}

func NewWebsocketEventService() *WebsocketEventService {
	return &WebsocketEventService{ps: pubsub.New(300)}
}

func (e *WebsocketEventService) Send(ctx context.Context, event Event) error {
	e.ps.Pub(event, string(event.Type))
	return nil
}

func (e *WebsocketEventService) Subscribe(t ...Type) chan interface{} {
	types := make([]string, len(t))
	for i, et := range t {
		types[i] = string(et)
	}
	return e.ps.Sub(types...)
}
