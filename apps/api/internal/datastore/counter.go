package datastore

import "context"

type Counter interface {
	Count(ctx context.Context, resource string) (int64, error)
}
