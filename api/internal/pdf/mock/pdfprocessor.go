package mock

import (
	"context"
	"github.com/dgmann/document-manager/api/internal/pdf"
	"github.com/dgmann/document-manager/api/internal/storage"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/stretchr/testify/mock"
	"io"
)

type PdfProcessor struct {
	mock.Mock
}

func NewPdfProcessor() *PdfProcessor {
	return &PdfProcessor{}
}

func (mock *PdfProcessor) Convert(ctx context.Context, f io.Reader, opts *pdf.ConvertOptions) ([]storage.Image, error) {
	args := mock.Called(ctx, f)
	return args.Get(0).([]storage.Image), args.Error(1)
}

func (mock *PdfProcessor) Rotate(ctx context.Context, image io.Reader, degrees int) (*storage.Image, error) {
	args := mock.Called(ctx, image, degrees)
	return args.Get(0).(*storage.Image), args.Error(1)
}

func (mock *PdfProcessor) Create(ctx context.Context, title string, records []api.Record) ([]byte, error) {
	args := mock.Called(ctx, title, records)
	return args.Get(0).([]byte), args.Error(1)
}
