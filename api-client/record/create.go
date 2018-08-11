package record

import (
	"io"
	"time"
)

type Repository struct {
	uploader *HttpUploader
}

func NewRepository(url string) *Repository {
	return &Repository{uploader: NewHttpUploader(url)}
}

type NewRecord struct {
	Sender       string
	ReceivedAt   time.Time
	File         io.Reader
	RetryCounter int
}

func (r *Repository) Create(record *NewRecord) error {
	return r.uploader.Upload(record)
}
