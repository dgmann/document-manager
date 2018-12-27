package http

import (
	"github.com/dgmann/document-manager/api/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"net/url"
	"time"
)

func NewTestContext() (*gin.Context, *gin.TestResponseRecorder) {
	w := gin.CreateTestResponseRecorder()
	context, _ := gin.CreateTestContext(w)
	context.Request = httptest.NewRequest("POST", "/", nil)
	requestUrl, _ := url.Parse("http://localhost")
	context.Set("location", requestUrl)
	return context, w
}

type TestFactory struct {
	*ResponseFactory
}

type mockModTimeReader struct {
	mock.Mock
}

func (m *mockModTimeReader) ModTime(resource app.KeyedResource) (time.Time, error) {
	args := m.Called(resource)
	return args.Get(0).(time.Time), args.Error(1)
}

func NewTestResponseFactory() (*TestFactory, *mockModTimeReader) {
	reader := new(mockModTimeReader)
	factory := NewResponseFactory(reader)
	return &TestFactory{factory}, reader
}
