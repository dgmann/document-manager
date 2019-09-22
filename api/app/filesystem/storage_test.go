package filesystem

import (
	"bytes"
	"errors"
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
	err := args.Error(1)
	res := args.Get(0)
	if res != nil {
		return res.(file), err
	}
	return nil, err
}

func (f filesystemMock) Open(name string) (file, error) {
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

func (f filesystemMock) Walk(p string, walkFn filepath.WalkFunc) error {
	args := f.Called(p, walkFn)
	return args.Error(0)
}

func (f *filesystemMock) locationExists(loc string) {
	if ext := filepath.Ext(loc); ext != "" {
		loc = filepath.Dir(loc)
	}
	f.On("Stat", filepath.FromSlash(loc)).Return(nil, nil)
}

type fileMock struct {
	mock.Mock
	*bytes.Buffer
}

func (f fileMock) Name() string {
	args := f.Called()
	return args.String(0)
}

func (f fileMock) Write(p []byte) (n int, err error) {
	_ = f.Called(p)
	return f.Buffer.Write(p)
}

func (f fileMock) Read(p []byte) (n int, err error) {
	_ = f.Called(p)
	return f.Buffer.Read(p)
}

func (f fileMock) Close() error {
	args := f.Called()
	return args.Error(0)
}

func (f *fileMock) mustBeClosed() {
	f.On("Close").Once().Return(nil)
}

func buildTestDiskStorage() (DiskStorage, *filesystemMock) {
	filesystem := new(filesystemMock)
	repository := DiskStorage{Root: filepath.FromSlash("/root"), storage: filesystem}
	return repository, filesystem
}

func TestRepository_Get(t *testing.T) {
	repository, filesystem := buildTestDiskStorage()
	resource := app.NewKeyedGenericResource([]byte{}, "png", "1")
	fileMock := new(fileMock)
	fileMock.mustBeClosed()
	fileMock.Buffer = bytes.NewBufferString("test")
	filename := filepath.FromSlash("/root/1.png")

	filesystem.On("Open", filename).Return(fileMock, nil)
	fileMock.On("Read", mock.Anything).Return()

	res, err := repository.Get(resource)
	assert.Nil(t, err)
	assert.Equal(t, []string{"1"}, res.Key())
	assert.Equal(t, "png", res.Format())
	assert.Equal(t, "test", string(res.Data()))

	filesystem.AssertExpectations(t)
	fileMock.AssertExpectations(t)
}

func TestRepository_Delete_File(t *testing.T) {
	repository, filesystem := buildTestDiskStorage()
	resource := app.NewKeyedGenericResource([]byte{}, "png", "1")

	filesystem.On("Remove", filepath.FromSlash("/root/1.png")).Return(nil)

	err := repository.Delete(resource)
	assert.Nil(t, err)

	filesystem.AssertExpectations(t)
}

func TestRepository_Delete_File_Multiple_Keys(t *testing.T) {
	repository, filesystem := buildTestDiskStorage()
	resource := app.NewKeyedGenericResource([]byte{}, "png", "1", "2")

	filesystem.On("Remove", filepath.FromSlash("/root/1/2.png")).Return(nil)

	err := repository.Delete(resource)
	assert.Nil(t, err)

	filesystem.AssertExpectations(t)
}

func TestRepository_Delete_Directory(t *testing.T) {
	repository, filesystem := buildTestDiskStorage()
	resource := app.NewKey("1")

	filesystem.On("RemoveAll", filepath.FromSlash("/root/1")).Return(nil)

	err := repository.Delete(resource)
	assert.Nil(t, err)

	filesystem.AssertExpectations(t)
}

func TestRepository_Normalize_Extension(t *testing.T) {
	repository, filesystem := buildTestDiskStorage()
	resource := app.NewKeyedGenericResource([]byte{}, "jpg", "1")

	filesystem.On("Remove", filepath.FromSlash("/root/1.jpeg")).Return(nil)

	err := repository.Delete(resource)
	assert.Nil(t, err)

	filesystem.AssertExpectations(t)
}

func TestRepository_Write(t *testing.T) {
	storage, filesystem := buildTestDiskStorage()
	resource := app.NewKeyedGenericResource([]byte{}, "png", "1", "2")
	fileMock := new(fileMock)
	fileMock.mustBeClosed()
	filename := "/root/1"

	filesystem.locationExists(filename)

	filesystem.On("Create", filepath.FromSlash(filename+"/2.png")).Return(fileMock, nil)
	fileMock.On("Write", []byte{}).Return(0, nil)

	err := storage.Write(resource)
	assert.Nil(t, err)

	filesystem.AssertExpectations(t)
	fileMock.AssertExpectations(t)
}

func TestRepository_Create_Fail(t *testing.T) {
	storage, filesystem := buildTestDiskStorage()
	resource := app.NewKeyedGenericResource([]byte{}, "png", "1", "2")
	filename := "/root/1/2.png"

	filesystem.locationExists(filename)
	filesystem.On("Create", filepath.FromSlash("/root/1/2.png")).Return(nil, errors.New(""))

	err := storage.Write(resource)
	assert.NotNil(t, err)

	filesystem.AssertExpectations(t)
}

func TestRepository_Write_Fail_Cleanup(t *testing.T) {
	storage, filesystem := buildTestDiskStorage()
	resource := app.NewKeyedGenericResource([]byte{}, "png", "1", "2")
	fileMock := new(fileMock)
	fileMock.mustBeClosed()
	filename := filepath.FromSlash("/root/1/2.png")

	filesystem.locationExists(filename)
	filesystem.On("Create", filename).Return(fileMock, nil)
	fileMock.On("Write", []byte{}).Return(0, errors.New(""))
	filesystem.On("Remove", filename).Once().Return(nil)

	err := storage.Write(resource)
	assert.Nil(t, err)

	filesystem.AssertExpectations(t)
	fileMock.AssertExpectations(t)
}
