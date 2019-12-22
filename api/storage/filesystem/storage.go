package filesystem

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/api/storage"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
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

//NewDiskStorage creates and initializes Storage on the local filesystem
// using the provided directory as the root.
func NewDiskStorage(directory string) (*DiskStorage, error) {
	storageEngine := osFileSystem{}
	diskstorage := &DiskStorage{storage: storageEngine, Root: directory}
	if err := diskstorage.ensureKeyedLocation(storage.NewKey(directory)); err != nil {
		return nil, fmt.Errorf("error creating disk storage: %w", err)
	}
	return diskstorage, nil
}

//Get retrieves a copy of the specified resource from the file system.
func (f *DiskStorage) Get(resource storage.Locatable) (*storage.GenericResource, error) {
	loc := f.Locate(resource)
	file, err := f.storage.Open(loc)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	resourceWithData := storage.NewKeyedGenericResource(data, resource.Format(), resource.Key()...)
	return resourceWithData, nil
}

//Delete deletes the specified Locatable from the file system.
func (f *DiskStorage) Delete(resource storage.Locatable) error {
	var err error
	loc := f.Locate(resource)
	if len(resource.Format()) > 0 {
		err = f.storage.Remove(loc)
	} else {
		err = f.storage.RemoveAll(loc)
	}
	if !os.IsNotExist(err) {
		return err
	}
	return nil
}

//Write writes the specified resource to the file system.
//It creates the files if it does not exist yet.
func (f *DiskStorage) Write(resource storage.KeyedResource) (err error) {
	if err := f.ensureResourceLocation(resource); err != nil {
		return err
	}

	loc := f.Locate(resource)
	file, err := f.storage.Create(loc)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", loc, err)
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = file.Write(resource.Data()); err != nil {
		return f.storage.Remove(loc)
	}

	return nil
}

type ForEachFunc func(resource storage.KeyedResource, err error) error

//ForEach executes the provided function for each stored element.
func (f *DiskStorage) ForEach(keyed storage.Keyed, forEachFn ForEachFunc) error {
	p := f.locate(keyed)
	return f.storage.Walk(p, func(currentPath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			logrus.Infof("Filename: %s", info.Name())
			ext := filepath.Ext(info.Name())
			abs := strings.Trim(currentPath, ext)
			rel, err := filepath.Rel(f.Root, abs)
			if err != nil {
				return err
			}

			keys := strings.Split(rel, string(filepath.Separator))
			resource := storage.NewKeyedGenericResource(nil, ext, keys...)
			logrus.Infof("Keys: %v", keys)

			withData, err := f.Get(resource)
			if err != nil {
				return fmt.Errorf("error reading resource %s: %w", filepath.Join(resource.Key()...), err)
			}
			return forEachFn(withData, err)
		}
		return nil
	})
}

//Check returns the health status of the file system.
func (f *DiskStorage) Check(ctx context.Context) (string, error) {
	if _, err := f.storage.Stat(f.Root); err != nil {
		return "", err
	}
	return "pass", nil
}

//ModTime returns the time at which the element was last changed.
func (f *DiskStorage) ModTime(resource storage.KeyedResource) (time.Time, error) {
	fp := f.Locate(resource)
	fileInfo, err := f.storage.Stat(fp)
	if err != nil {
		return time.Now(), err
	}
	return fileInfo.ModTime(), nil
}

//LocateResource returns the location of the specified Locatable.
func (f *DiskStorage) Locate(resource storage.Locatable) string {
	dir := f.locate(resource)
	if len(resource.Format()) > 0 {
		format := normalizeExtension(resource.Format())
		return fmt.Sprintf("%s.%s", dir, format)
	}
	return dir
}

//Locate returns the location of the specified Keyed element
func (f *DiskStorage) locate(keyed storage.Keyed) string {
	keySlice := append([]string{f.Root}, keyed.Key()...)
	return filepath.Join(keySlice...)
}

func normalizeExtension(extension string) string {
	if extension == "jpg" {
		return "jpeg"
	}
	return extension
}

//ensureResourceLocation ensures the the directory structure required to store the KeyedResource is in place.
func (f *DiskStorage) ensureResourceLocation(keyed storage.KeyedResource) error {
	loc := f.Locate(keyed)
	dir := filepath.Dir(loc)
	return f.ensureLocation(dir)
}

//ensureResourceLocation ensures the the directory structure required to store the Keyed is in place.
func (f *DiskStorage) ensureKeyedLocation(keyed storage.Keyed) error {
	dir := f.locate(keyed)
	return f.ensureLocation(dir)
}

func (f *DiskStorage) ensureLocation(dir string) error {
	if _, err := f.storage.Stat(dir); os.IsNotExist(err) {
		err = f.storage.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create directory %s: %w", dir, err)
		}
	}
	return nil
}
