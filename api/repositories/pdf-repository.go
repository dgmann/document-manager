package repositories

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type PDFRepository interface {
	Get(id string) (io.Reader, error)
	ResourceWriter
}

type FileSystemPDFRepository struct {
	*FileSystemRepository
	directory string
}

func NewFileSystemPDFRepository(directory string) PDFRepository {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, os.ModePerm)
	}
	return &FileSystemPDFRepository{directory: directory, FileSystemRepository: NewFileSystemRepository(directory)}
}

func (f *FileSystemPDFRepository) Get(id string) (io.Reader, error) {
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

func (f *FileSystemPDFRepository) Set(resource KeyedResource) error {
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

func (f *FileSystemPDFRepository) Delete(resource KeyedResource) error {
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
