package filesystem

import (
	"bytes"
	"github.com/dgmann/document-manager/api/app"
	"io"
	"io/ioutil"
	"os"
	"path"
)

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
	fp := path.Join(f.baseDirectory, id+".pdf")
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

func (f *ArchiveService) Set(resource app.KeyedResource) error {
	keys := append([]string{f.baseDirectory}, resource.Key()...)
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
