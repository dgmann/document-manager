package repositories

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
)

type FileSystemRepository struct {
	baseDirectory string
}

func NewFileSystemRepository(baseDirectory string) *FileSystemRepository {
	return &FileSystemRepository{baseDirectory: baseDirectory}
}

func (f *FileSystemRepository) Delete(resource KeyedResource) error {
	p := f.buildPath(resource.Key()...)
	var err error
	if len(resource.Format()) > 0 {
		p += "." + normalizeExtension(resource.Format())
		err = os.Remove(p)
	} else {
		err = os.RemoveAll(p)
	}
	if !os.IsNotExist(err) {
		return err
	}
	logrus.Infof("%s cannot be deleted as it does not exist", p)
	return nil
}

func (f *FileSystemRepository) Write(resource KeyedResource) (err error) {
	fp := f.buildResourcePath(resource)

	dir := filepath.Dir(fp)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	imageFile, err := os.Create(fp)
	defer imageFile.Close()
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovering: %v", r)
			os.Remove(fp)
			err = errors.New("failed to save images")
		}
	}()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"name":      imageFile.Name(),
			"directory": dir,
			"error":     err,
		}).Error("Error creating image file")
		return err
	}
	_, err = imageFile.Write(resource.Data())
	return err
}

func (f *FileSystemRepository) buildResourcePath(resource KeyedResource) string {
	p := f.buildPath(resource.Key()...)
	if len(resource.Format()) > 0 {
		p += "." + normalizeExtension(resource.Format())
	}
	return p
}

func (f *FileSystemRepository) buildPath(keys ...string) string {
	keySlice := append([]string{f.baseDirectory}, keys...)
	return path.Join(keySlice...)
}

func normalizeExtension(extension string) string {
	if extension == "jpeg" {
		return "jpg"
	}
	return extension
}
