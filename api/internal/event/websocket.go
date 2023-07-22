package event

import (
	"context"
	"github.com/cskr/pubsub"
	"github.com/dgmann/document-manager/api/pkg/api"
)

type WebsocketEventService[T any] struct {
	ps *pubsub.PubSub
}

func NewWebsocketEventService[T any]() *WebsocketEventService[T] {
	return &WebsocketEventService[T]{ps: pubsub.New(300)}
}

func (e *WebsocketEventService[T]) Send(ctx context.Context, event api.Event[T]) error {
	e.ps.Pub(event, string(event.Type))
	return nil
}

func (e *WebsocketEventService[T]) Subscribe(t ...api.EventType) chan interface{} {
	types := make([]string, len(t))
	for i, et := range t {
		types[i] = string(et)
	}
	return e.ps.Sub(types...)
}
