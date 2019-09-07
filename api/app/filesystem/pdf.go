package filesystem

import (
	"bytes"
	"github.com/dgmann/document-manager/api/app"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const PDFFileExtension = ".pdf"

type ArchiveService struct {
	*Storage
}

func NewArchiveService(directory string) (*ArchiveService, error) {
	repository, err := New(directory)
	if err != nil {
		return nil, err
	}
	return &ArchiveService{Storage: repository}, nil
}

func (f *ArchiveService) Get(id string) (io.Reader, error) {
	fp := path.Join(f.baseDirectory, id+PDFFileExtension)
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

func (f *ArchiveService) Set(resource app.KeyedResource) (err error) {
	keys := append([]string{f.baseDirectory}, resource.Key()...)
	p := path.Join(keys...)

	fp := p + "." + resource.Format()
	pdfFile, err := os.Create(fp)
	if err != nil {
		return
	}
	defer func() {
		cerr := pdfFile.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = pdfFile.Write(resource.Data())
	return
}

func (f *Storage) NumberOfElements() (int, error) {
	count := 0
	err := filepath.Walk(f.baseDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == PDFFileExtension {
				count++
			}
			return nil
		})
	if err != nil {
		return count, err
	}
	return count, nil
}
