package http

import (
	"bytes"
	"encoding/json"
	"github.com/dgmann/document-manager/api/app"
	"github.com/dgmann/document-manager/api/app/mock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
)

func createTestController() (*CategoryController, *mock.CategoryService) {
	mockCategoryRepository := new(mock.CategoryService)
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
