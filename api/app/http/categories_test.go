package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dgmann/document-manager/api/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) All(ctx context.Context) ([]app.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]app.Category), args.Error(1)
}

func (m *MockCategoryRepository) Add(ctx context.Context, id, category string) error {
	args := m.Called(ctx, id, category)
	return args.Error(0)
}

func createTestController() (*CategoryController, *MockCategoryRepository) {
	mockCategoryRepository := new(MockCategoryRepository)
	testResponseService, _ := NewTestResponseFactory()
	return &CategoryController{categories: mockCategoryRepository, responseService: testResponseService}, mockCategoryRepository
}

func TestCategoryController_All(t *testing.T) {
	controller, mockCategoryRepository := createTestController()

	ctx, w := NewTestContext()
	categories := []app.Category{{Name: "mock", Id: "mock"}}
	mockCategoryRepository.On("All", ctx).Return(categories, nil)

	controller.All(ctx)
	expected, _ := json.Marshal(categories)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
}

func TestCategoryController_Create(t *testing.T) {
	controller, mockCategoryRepository := createTestController()

	ctx, w := NewTestContext()

	category := app.Category{Id: "cat", Name: "Category"}
	body, _ := json.Marshal(category)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	ctx.Request = req

	mockCategoryRepository.On("Add", ctx, category.Id, category.Name).Return(nil)

	controller.Create(ctx)
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, string(body), strings.TrimSpace(w.Body.String()))
}
