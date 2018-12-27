package filesystem

import (
	"io"
	"os"
)

type filesystem interface {
	Remove(name string) error
	RemoveAll(path string) error
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (file, error)
	Stat(name string) (os.FileInfo, error)
}

type file interface {
	io.WriteCloser
	Name() string
}

type diskFileSystem struct {
}

func (f diskFileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (f diskFileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (f diskFileSystem) Create(name string) (file, error) {
	return os.Create(name)
}

func (f diskFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f diskFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
