package filesystem

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/api/internal/storage"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type fileSystem interface {
	Remove(name string) error
	RemoveAll(name string) error
	Create(name string) (file, error)
	Open(name string) (file, error)
	MkdirAll(path string, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	Walk(p string, walkFn filepath.WalkFunc) error
}

type file interface {
	io.Closer
	io.Reader
	io.Writer
}

type DiskStorage struct {
	Root    string
	storage fileSystem
}

// NewDiskStorage creates and initializes Storage on the local filesystem
// using the provided directory as the root.
func NewDiskStorage(directory string) (*DiskStorage, error) {
	storageEngine := osFileSystem{}
	diskstorage := &DiskStorage{storage: storageEngine, Root: directory}
	if err := diskstorage.ensureKeyedLocation(storage.RootKey); err != nil {
		return nil, fmt.Errorf("error creating disk storage: %w", err)
	}
	return diskstorage, nil
}

// Get retrieves a copy of the specified resource from the file system.
func (s *DiskStorage) Get(resource storage.Locatable) (resourceWithData *storage.GenericResource, err error) {
	loc := s.Locate(resource)
	openedFile, err := s.storage.Open(loc)
	if err != nil {
		return nil, err
	}
	defer func(f file) {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}(openedFile)

	data, err := io.ReadAll(openedFile)
	if err != nil {
		return nil, err
	}
	resourceWithData = storage.NewKeyedGenericResource(data, resource.Format(), resource.Key()...)
	return resourceWithData, nil
}

// Delete deletes the specified Locatable from the file system.
func (s *DiskStorage) Delete(resource storage.Locatable) error {
	var err error
	loc := s.Locate(resource)
	if len(resource.Format()) > 0 {
		err = s.storage.Remove(loc)
	} else {
		err = s.storage.RemoveAll(loc)
	}
	if !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Write writes the specified resource to the file system.
// It creates the files if it does not exist yet.
func (s *DiskStorage) Write(resource storage.KeyedResource) (err error) {
	if err := s.ensureResourceLocation(resource); err != nil {
		return err
	}

	loc := s.Locate(resource)
	newFile, err := s.storage.Create(loc)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", loc, err)
	}
	defer func(f file) {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}(newFile)

	if _, err = newFile.Write(resource.Data()); err != nil {
		return s.storage.Remove(loc)
	}

	return nil
}

type ForEachFunc func(resource storage.KeyedResource, err error) error

// ForEach executes the provided function for each stored element.
func (s *DiskStorage) ForEach(keyed storage.Keyed, forEachFn ForEachFunc) error {
	p := s.locate(keyed)
	return s.storage.Walk(p, func(currentPath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			ext := filepath.Ext(info.Name())
			abs := strings.TrimSuffix(currentPath, ext)
			rel, err := filepath.Rel(s.Root, abs)
			if err != nil {
				return err
			}

			keys := strings.Split(rel, string(filepath.Separator))
			resource := storage.NewKeyedGenericResource(nil, ext, keys...)
			withData, err := s.Get(resource)
			if err != nil {
				return fmt.Errorf("error reading resource %s: %w", filepath.Join(resource.Key()...), err)
			}
			return forEachFn(withData, err)
		}
		return nil
	})
}

// Check returns the health status of the file system.
func (s *DiskStorage) Check(ctx context.Context) (string, error) {
	if _, err := s.storage.Stat(s.Root); err != nil {
		return "", err
	}
	return "pass", nil
}

// ModTime returns the time at which the element was last changed.
func (s *DiskStorage) ModTime(resource storage.KeyedResource) (time.Time, error) {
	fp := s.Locate(resource)
	fileInfo, err := s.storage.Stat(fp)
	if err != nil {
		return time.Now(), err
	}
	return fileInfo.ModTime(), nil
}

// Locate returns the location of the specified Locatable.
func (s *DiskStorage) Locate(resource storage.Locatable) string {
	dir := s.locate(resource)
	if len(resource.Format()) > 0 {
		format := normalizeExtension(resource.Format())
		return fmt.Sprintf("%s.%s", dir, format)
	}
	return dir
}

// locate returns the location of the specified Keyed element.
func (s *DiskStorage) locate(keyed storage.Keyed) string {
	keySlice := append([]string{s.Root}, keyed.Key()...)
	return filepath.Join(keySlice...)
}

func normalizeExtension(extension string) string {
	if extension == "jpg" {
		return "jpeg"
	}
	return extension
}

// ensureResourceLocation ensures the directory structure required to store the KeyedResource is in place.
func (s *DiskStorage) ensureResourceLocation(keyed storage.KeyedResource) error {
	loc := s.Locate(keyed)
	dir := filepath.Dir(loc)
	return s.ensureLocation(dir)
}

// ensureResourceLocation ensures the directory structure required to store the Keyed is in place.
func (s *DiskStorage) ensureKeyedLocation(keyed storage.Keyed) error {
	dir := s.locate(keyed)
	return s.ensureLocation(dir)
}

func (s *DiskStorage) ensureLocation(dir string) error {
	if _, err := s.storage.Stat(dir); os.IsNotExist(err) {
		err = s.storage.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create directory %s: %w", dir, err)
		}
	}
	return nil
}
