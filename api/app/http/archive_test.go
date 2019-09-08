package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"io"
	"net/http/httptest"
	"testing"
)

type MockPdfRepository struct {
	mock.Mock
}

func (m *MockPdfRepository) Get(id string) (io.Reader, error) {
	args := m.Called(id)
	return args.Get(0).(io.Reader), args.Error(1)
}

func TestArchiveController_One(t *testing.T) {
	mockPdfRepository := new(MockPdfRepository)
	controller := ArchiveController{pdfs: mockPdfRepository}

	mockFile := bytes.NewBufferString("mock")
	mockPdfRepository.On("Get", "1").Return(mockFile, nil)

	req := httptest.NewRequest("Get", "/1", nil)
	ctx := req.Context()
	req = req.WithContext(context.WithValue(ctx, "recordId", "1"))
	w := httptest.NewRecorder()
	controller.One(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "mock", w.Body.String())
}
