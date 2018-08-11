package record

import (
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
	PdfPath      string
	RetryCounter int
}

func (r *Repository) Create(record *NewRecord) {
	r.uploader.Upload(record)
}
