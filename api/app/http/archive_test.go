package http

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
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

	w := gin.CreateTestResponseRecorder()
	context, _ := gin.CreateTestContext(w)
	context.Params = gin.Params{gin.Param{Key: "recordId", Value: "1"}}
	mockFile := bytes.NewBufferString("mock")
	mockPdfRepository.On("Get", "1").Return(mockFile, nil)

	controller.One(context)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "mock", w.Body.String())
}
