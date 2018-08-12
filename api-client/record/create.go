package record

import (
	"io"
	"github.com/dgmann/document-manager/api/models"
)

type Repository struct {
	uploader *HttpUploader
}

func NewRepository(url string) *Repository {
	return &Repository{uploader: NewHttpUploader(url)}
}

type CreateRecord = models.CreateRecord
type Status = models.Status

const (
	StatusInbox     = models.StatusInbox
	StatusEscalated = models.StatusEscalated
	StatusReview    = models.StatusReview
	StatusOther     = models.StatusOther
	StatusDone      = models.StatusDone
)

type NewRecord struct {
	CreateRecord
	File         io.Reader
	RetryCounter int
}

func (r *Repository) Create(record *NewRecord) error {
	return r.uploader.Upload(record)
}
