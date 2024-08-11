package datastore

import "context"

type TagService interface {
	All(ctx context.Context) ([]string, error)
	ByPatient(ctx context.Context, id string) ([]string, error)
}
