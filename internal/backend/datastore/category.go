package datastore

import (
	"context"

	"github.com/dgmann/document-manager/pkg/api"
)

type CategoryService interface {
	All(ctx context.Context) ([]api.Category, error)
	Find(ctx context.Context, id string) (*api.Category, error)
	FindByPatient(ctx context.Context, id string) ([]api.Category, error)
	Add(ctx context.Context, category *api.Category) error
	Update(ctx context.Context, category *api.Category) error
}
