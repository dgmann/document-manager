package filesystem

import (
	"context"
	"github.com/dgmann/document-manager/api/app"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	baseDirectory string
	filesystem    filesystem
}

func New(directory string) (*Storage, error) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if e := os.MkdirAll(directory, os.ModePerm); e != nil {
			logrus.WithError(e).Error("error creating directory")
		}
	}
	return &Storage{baseDirectory: directory, filesystem: diskFileSystem{}}, nil
}

func (f *Storage) Delete(resource app.KeyedResource) error {
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

func (f *Storage) Write(resource app.KeyedResource) (err error) {
	fp := f.buildResourcePath(resource)

	dir := filepath.Dir(fp)
	if _, err := f.filesystem.Stat(dir); os.IsNotExist(err) {
		err = f.filesystem.MkdirAll(dir, os.ModePerm)
		logrus.WithError(err).Error("could not create directory")
	}

	imageFile, err := f.filesystem.Create(fp)
	defer imageFile.Close()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"name":      imageFile.Name(),
			"directory": dir,
			"error":     err,
		}).Error("Error creating file")
		return err
	}
	_, err = imageFile.Write(resource.Data())
	if err != nil {
		logrus.WithError(err).Error("error writing file. Cleanup leftovers")
		return f.filesystem.Remove(fp)
	}
	return nil
}

func (f *Storage) Check(ctx context.Context) (string, error) {
	if _, err := os.Stat(f.baseDirectory); err != nil {
		return "", err
	}
	return "pass", nil
}

func (f *Storage) ModTime(resource app.KeyedResource) (time.Time, error) {
	fp := f.buildResourcePath(resource)
	fileInfo, err := os.Stat(fp)
	if err != nil {
		return time.Now(), err
	}
	return fileInfo.ModTime(), nil
}

func (f *Storage) buildResourcePath(resource app.KeyedResource) string {
	p := f.buildPath(resource.Key()...)
	if len(resource.Format()) > 0 {
		p += "." + normalizeExtension(resource.Format())
	}
	return p
}

func (f *Storage) buildPath(keys ...string) string {
	keySlice := append([]string{f.baseDirectory}, keys...)
	return filepath.Join(keySlice...)
}

func normalizeExtension(extension string) string {
	if extension == "jpg" {
		return "jpeg"
	}
	return extension
}
