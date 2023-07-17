package filesystem

import (
	"github.com/dgmann/document-manager/api/internal/storage"
)

const PDFFileExtension = "pdf"

type ArchiveService struct {
	*DiskStorage
}

func NewArchiveService(directory string) (*ArchiveService, error) {
	repository, err := NewDiskStorage(directory)
	if err != nil {
		return nil, err
	}
	return &ArchiveService{DiskStorage: repository}, nil
}

type RecordDirectory struct {
	recordId string
}

func (dir RecordDirectory) Key() []string {
	return []string{dir.recordId}
}

func (f *ArchiveService) Get(id string) (storage.KeyedResource, error) {
	resource := storage.NewKeyedGenericResource(nil, PDFFileExtension, id)
	return f.DiskStorage.Get(resource)
}

func (f *ArchiveService) Set(resource storage.KeyedResource) error {
	return f.Write(resource)
}

func (s *DiskStorage) NumberOfElements() (int, error) {
	count := 0
	err := s.ForEach(storage.NewKey("/"), func(resource storage.KeyedResource, err error) error {
		if err != nil {
			return err
		}
		if resource.Format() == PDFFileExtension {
			count++
		}
		return nil
	})
	if err != nil {
		return count, err
	}
	return count, nil
}
