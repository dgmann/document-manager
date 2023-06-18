package event

import (
	"context"
	"github.com/dgmann/document-manager/api/pkg/api"
)

type Subscriber interface {
	Subscribe(t ...api.Type) chan interface{}
}

type Sender interface {
	Send(ctx context.Context, event api.Event) error
}
