package event

import (
	"context"
	"github.com/cskr/pubsub"
	"github.com/dgmann/document-manager/api/pkg/api"
)

type WebsocketEventService struct {
	ps *pubsub.PubSub
}

func NewWebsocketEventService() *WebsocketEventService {
	return &WebsocketEventService{ps: pubsub.New(300)}
}

func (e *WebsocketEventService) Send(ctx context.Context, event api.Event) error {
	e.ps.Pub(event, string(event.Type))
	return nil
}

func (e *WebsocketEventService) Subscribe(t ...api.Type) chan interface{} {
	types := make([]string, len(t))
	for i, et := range t {
		types[i] = string(et)
	}
	return e.ps.Sub(types...)
}
