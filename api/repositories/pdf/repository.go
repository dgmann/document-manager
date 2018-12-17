package pdf

import (
	"bytes"
	"github.com/dgmann/document-manager/api/repositories"
	"github.com/dgmann/document-manager/api/repositories/filesystem"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type Repository interface {
	Get(id string) (io.Reader, error)
	repositories.ResourceWriter
}

type FileSystemRepository struct {
	*filesystem.Repository
	directory string
}

func NewFileSystemRepository(directory string) Repository {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, os.ModePerm)
	}
	return &FileSystemRepository{directory: directory, Repository: filesystem.NewRepository(directory)}
}

func (f *FileSystemRepository) Get(id string) (io.Reader, error) {
	fp := path.Join(f.directory, id+".pdf")
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

func (f *FileSystemRepository) Set(resource repositories.KeyedResource) error {
	keys := append([]string{f.directory}, resource.Key()...)
	p := path.Join(keys...)

	fp := p + "." + resource.Format()
	pdfFile, err := os.Create(fp)
	defer pdfFile.Close()
	if err != nil {
		return err
	}
	_, err = pdfFile.Write(resource.Data())
	return err
}

func (f *FileSystemRepository) Delete(resource repositories.KeyedResource) error {
	keys := append([]string{f.directory}, resource.Key()...)
	p := path.Join(keys...)

	fp := p + "." + resource.Format()

	err := os.Remove(fp)
	if !os.IsNotExist(err) {
		return err
	}
	log.Infof("%s cannot be deleted as it does not exist", fp)
	return nil
}
