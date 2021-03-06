package mock

import (
	"context"
	"github.com/dgmann/document-manager/api/datastore"
	"github.com/stretchr/testify/mock"
)

type CategoryService struct {
	mock.Mock
}

func (m *CategoryService) All(ctx context.Context) ([]datastore.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]datastore.Category), args.Error(1)
}

func (m *CategoryService) Add(ctx context.Context, id, category string) error {
	args := m.Called(ctx, id, category)
	return args.Error(0)
}
