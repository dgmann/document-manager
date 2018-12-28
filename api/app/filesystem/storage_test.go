package filesystem

import (
	"github.com/dgmann/document-manager/api/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"path/filepath"
	"testing"
)

type filesystemMock struct {
	mock.Mock
}

func (f filesystemMock) Remove(name string) error {
	args := f.Called(name)
	return args.Error(0)
}

func (f filesystemMock) RemoveAll(path string) error {
	args := f.Called(path)
	return args.Error(0)
}

func (f filesystemMock) Create(name string) (file, error) {
	args := f.Called(name)
	return args.Get(0).(file), args.Error(1)
}

func (f filesystemMock) MkdirAll(path string, perm os.FileMode) error {
	args := f.Called(path)
	return args.Error(0)
}

func (f filesystemMock) Stat(name string) (os.FileInfo, error) {
	args := f.Called(name)
	return nil, args.Error(1)
}

type fileMock struct {
	mock.Mock
}

func (f fileMock) Name() string {
	args := f.Called()
	return args.String(0)
}

func (f fileMock) Write(p []byte) (n int, err error) {
	args := f.Called(p)
	return args.Int(0), args.Error(1)
}

func (f fileMock) Close() error {
	args := f.Called()
	return args.Error(0)
}

func buildTestRepository() (Storage, *filesystemMock) {
	filesystem := new(filesystemMock)
	repository := Storage{baseDirectory: filepath.FromSlash("/root"), filesystem: filesystem}
	return repository, filesystem
}

func TestRepository_Delete_File(t *testing.T) {
	repository, filesystem := buildTestRepository()
	resource := app.NewKeyedGenericResource([]byte{}, "png", "1")

	filesystem.On("Remove", filepath.FromSlash("/root/1.png")).Return(nil)

	err := repository.Delete(resource)
	assert.Nil(t, err)

	filesystem.AssertExpectations(t)
}

func TestRepository_Delete_File_Multiple_Keys(t *testing.T) {
	repository, filesystem := buildTestRepository()
	resource := app.NewKeyedGenericResource([]byte{}, "png", "1", "2")

	filesystem.On("Remove", filepath.FromSlash("/root/1/2.png")).Return(nil)

	err := repository.Delete(resource)
	assert.Nil(t, err)

	filesystem.AssertExpectations(t)
}

func TestRepository_Delete_Directory(t *testing.T) {
	repository, filesystem := buildTestRepository()
	resource := app.NewDirectoryResource("1")

	filesystem.On("RemoveAll", filepath.FromSlash("/root/1")).Return(nil)

	err := repository.Delete(resource)
	assert.Nil(t, err)

	filesystem.AssertExpectations(t)
}

func TestRepository_Normalize_Extension(t *testing.T) {
	repository, filesystem := buildTestRepository()
	resource := app.NewKeyedGenericResource([]byte{}, "jpg", "1")

	filesystem.On("Remove", filepath.FromSlash("/root/1.jpeg")).Return(nil)

	err := repository.Delete(resource)
	assert.Nil(t, err)

	filesystem.AssertExpectations(t)
}

func TestRepository_Write(t *testing.T) {
	repository, filesystem := buildTestRepository()
	resource := app.NewKeyedGenericResource([]byte{}, "png", "1", "2")
	fileMock := new(fileMock)

	filesystem.On("Stat", filepath.FromSlash("/root/1")).Return(nil, os.ErrNotExist)
	filesystem.On("MkdirAll", filepath.FromSlash("/root/1")).Once().Return(nil)
	filesystem.On("Create", filepath.FromSlash("/root/1/2.png")).Return(fileMock, nil)
	fileMock.On("Close").Once().Return(nil)
	fileMock.On("Name").Return("Mock")
	fileMock.On("Write", []byte{}).Return(0, nil)

	err := repository.Write(resource)
	assert.Nil(t, err)

	filesystem.AssertExpectations(t)
}
