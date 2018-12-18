package http

import (
	"bytes"
	"encoding/json"
	"github.com/dgmann/document-manager/api/http/response"
	"github.com/dgmann/document-manager/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"strings"
	"testing"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) All() ([]models.Category, error) {
	args := m.Called()
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoryRepository) Add(id, category string) error {
	args := m.Called(id, category)
	return args.Error(0)
}

func createTestController() (*CategoryController, *MockCategoryRepository) {
	mockCategoryRepository := new(MockCategoryRepository)
	testResponseService, _ := response.NewTestFactory()
	return &CategoryController{categories: mockCategoryRepository, responseService: testResponseService}, mockCategoryRepository
}

func TestCategoryController_All(t *testing.T) {
	controller, mockCategoryRepository := createTestController()

	context, w := NewTestContext()
	categories := []models.Category{{Name: "mock", Id: "mock"}}
	mockCategoryRepository.On("All").Return(categories, nil)

	controller.All(context)
	expected, _ := json.Marshal(categories)
	assert.Equal(t, string(expected), strings.TrimSpace(w.Body.String()))
}

func TestCategoryController_Create(t *testing.T) {
	controller, mockCategoryRepository := createTestController()

	context, w := NewTestContext()

	category := models.Category{Id: "cat", Name: "Category"}
	body, _ := json.Marshal(category)
	req, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	context.Request = req

	mockCategoryRepository.On("Add", category.Id, category.Name).Return(nil)

	controller.Create(context)
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, string(body), strings.TrimSpace(w.Body.String()))
}
