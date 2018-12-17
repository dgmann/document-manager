package filesystem

import (
	"errors"
	"github.com/dgmann/document-manager/api/repositories"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Repository struct {
	baseDirectory string
	filesystem    filesystem
}

func NewRepository(baseDirectory string) *Repository {
	return &Repository{baseDirectory: baseDirectory, filesystem: diskFileSystem{}}
}

func (f *Repository) Delete(resource repositories.KeyedResource) error {
	p := f.buildPath(resource.Key()...)
	var err error
	if len(resource.Format()) > 0 {
		p += "." + normalizeExtension(resource.Format())
		err = f.filesystem.Remove(p)
	} else {
		err = f.filesystem.RemoveAll(p)
	}
	if !os.IsNotExist(err) {
		return err
	}
	logrus.Infof("%s cannot be deleted as it does not exist", p)
	return nil
}

func (f *Repository) Write(resource repositories.KeyedResource) (err error) {
	fp := f.buildResourcePath(resource)

	dir := filepath.Dir(fp)
	if _, err := f.filesystem.Stat(dir); os.IsNotExist(err) {
		err = f.filesystem.MkdirAll(dir, os.ModePerm)
		logrus.WithError(err).Error("could not create directory")
	}

	imageFile, err := f.filesystem.Create(fp)
	defer imageFile.Close()
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovering: %v", r)
			f.filesystem.Remove(fp)
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

func (f *Repository) buildResourcePath(resource repositories.KeyedResource) string {
	p := f.buildPath(resource.Key()...)
	if len(resource.Format()) > 0 {
		p += "." + normalizeExtension(resource.Format())
	}
	return p
}

func (f *Repository) buildPath(keys ...string) string {
	keySlice := append([]string{f.baseDirectory}, keys...)
	return filepath.Join(keySlice...)
}

func normalizeExtension(extension string) string {
	if extension == "jpeg" {
		return "jpg"
	}
	return extension
}
