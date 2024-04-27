package filesystem

import (
	"os"
	"path/filepath"
)

type osFileSystem struct {
}

func (f osFileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (f osFileSystem) RemoveAll(name string) error {
	return os.RemoveAll(name)
}

func (f osFileSystem) Create(name string) (file, error) {
	return os.Create(name)
}

func (f osFileSystem) Open(name string) (file, error) {
	return os.Open(name)
}

func (f osFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f osFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (f osFileSystem) Walk(p string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(p, walkFn)
}
