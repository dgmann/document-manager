package http

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgmann/document-manager/internal/backend/datastore/mock"
	"github.com/dgmann/document-manager/pkg/api"
	"github.com/stretchr/testify/assert"
)

func createTestController() (*CategoryController, *mock.CategoryService) {
	mockCategoryRepository := new(mock.CategoryService)
	return &CategoryController{categories: mockCategoryRepository}, mockCategoryRepository
}

func TestCategoryController_All(t *testing.T) {
	controller, mockCategoryRepository := createTestController()

	categories := []api.Category{{Name: "mock", Id: "mock"}}
	mockCategoryRepository.On("All", mock.Anything).Return(categories, nil)

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	controller.All(w, req)
	expected, _ := json.Marshal(categories)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
}

func TestCategoryController_Create(t *testing.T) {
	controller, mockCategoryRepository := createTestController()

	category := &api.Category{Id: "cat", Name: "Category"}
	mockCategoryRepository.On("Add", mock.Anything, category).Return(nil)

	body, _ := json.Marshal(category)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	controller.Create(w, req)
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, string(body), strings.TrimSpace(w.Body.String()))
}
