package event

import (
	"context"

	"github.com/dgmann/document-manager/pkg/api"
)

type Subscriber interface {
	Subscribe(t ...api.EventType) chan interface{}
}

type Sender[T any] interface {
	Send(ctx context.Context, event api.Event[T]) error
}
