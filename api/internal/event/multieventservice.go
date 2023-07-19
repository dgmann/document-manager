package event

import (
	"context"
	"errors"
	"github.com/dgmann/document-manager/api/pkg/api"
)

type MultiEventSender[T any] struct {
	senders []Sender[T]
}

func NewMultiEventSender[T any](senders ...Sender[T]) *MultiEventSender[T] {
	return &MultiEventSender[T]{senders: senders}
}

func (e *MultiEventSender[T]) Send(ctx context.Context, event api.Event[T]) error {
	var errs []error
	for _, sender := range e.senders {
		if err := sender.Send(ctx, event); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
