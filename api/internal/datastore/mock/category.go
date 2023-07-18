package mock

import (
	"context"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/stretchr/testify/mock"
)

type CategoryService struct {
	mock.Mock
}

func (m *CategoryService) All(ctx context.Context) ([]api.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]api.Category), args.Error(1)
}

func (m *CategoryService) Add(ctx context.Context, category *api.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *CategoryService) Update(ctx context.Context, category *api.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}
