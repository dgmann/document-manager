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
	Set(id string, data []byte) error
	Delete(id string) error
}

type FileSystemPDFRepository struct {
	directory string
}

func newFileSystemPDFRepository(directory string) PDFRepository {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, os.ModePerm)
	}
	return &FileSystemPDFRepository{directory: directory}
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

func (f *FileSystemPDFRepository) Set(id string, data []byte) error {
	filePath := path.Join(f.directory, id+".pdf")
	pdfFile, err := os.Create(filePath)
	defer pdfFile.Close()
	if err != nil {
		return err
	}
	_, err = pdfFile.Write(data)
	return err
}

func (f *FileSystemPDFRepository) Delete(id string) error {
	fp := path.Join(f.directory, id+".pdf")
	err := os.Remove(fp)
	if !os.IsNotExist(err) {
		return err
	}
	log.Infof("%s cannot be deleted as it does not exist", fp)
	return nil
}
