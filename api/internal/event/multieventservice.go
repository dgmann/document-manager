package event

import (
	"context"
	"errors"
	"github.com/dgmann/document-manager/api/pkg/api"
)

type MultiEventSender struct {
	senders []Sender
}

func NewMultiEventSender(senders ...Sender) *MultiEventSender {
	return &MultiEventSender{senders: senders}
}

func (e *MultiEventSender) Send(ctx context.Context, event api.Event) error {
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
