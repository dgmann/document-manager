package response

import (
	"github.com/dgmann/document-manager/api/repositories"
	"github.com/stretchr/testify/mock"
	"time"
)

type TestFactory struct {
	*Factory
}

type mockModTimeReader struct {
	mock.Mock
}

func (m *mockModTimeReader) ModTime(resource repositories.KeyedResource) (time.Time, error) {
	args := m.Called(resource)
	return args.Get(0).(time.Time), args.Error(1)
}

func NewTestFactory() (*TestFactory, *mockModTimeReader) {
	reader := new(mockModTimeReader)
	factory := NewFactory(reader)
	return &TestFactory{factory}, reader
}
