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
	return &CategoryController{categories: mockCategoryRepository}, mockCategoryRepository
}

func TestCategoryController_All(t *testing.T) {
	controller, mockCategoryRepository := createTestController()

	categories := []app.Category{{Name: "mock", Id: "mock"}}
	mockCategoryRepository.On("All", mock.Anything).Return(categories, nil)

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	controller.All(w, req)
	expected, _ := json.Marshal(categories)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
}

func TestCategoryController_Create(t *testing.T) {
	controller, mockCategoryRepository := createTestController()

	category := app.Category{Id: "cat", Name: "Category"}
	mockCategoryRepository.On("Add", mock.Anything, category.Id, category.Name).Return(nil)

	body, _ := json.Marshal(category)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	controller.Create(w, req)
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, string(body), strings.TrimSpace(w.Body.String()))
}
