package client

import (
	"github.com/dgmann/document-manager/api/datastore"
	"io"
)

type Repository struct {
	uploader *HttpUploader
}

func NewRepository(url string) *Repository {
	return &Repository{uploader: NewHttpUploader(url)}
}

type NewRecord struct {
	datastore.CreateRecord
	File         io.Reader
	RetryCounter int
}

func (r *Repository) Create(record *NewRecord) error {
	return r.uploader.Upload(record)
}
