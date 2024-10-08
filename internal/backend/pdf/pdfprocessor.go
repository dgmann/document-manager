package pdf

import (
	"context"
	"io"

	"github.com/dgmann/document-manager/internal/backend/storage"
	"github.com/dgmann/document-manager/pkg/api"
)

type Processor interface {
	Converter
	Rotator
	Creator
}

type Method string

func (m Method) String() string {
	return string(m)
}

const (
	EXTRACT   Method = "extract"
	RASTERIZE Method = "rasterize"
)

type ConvertOptions struct {
	Method Method
}

func Extract() *ConvertOptions {
	return &ConvertOptions{Method: EXTRACT}
}

func Rasterize() *ConvertOptions {
	return &ConvertOptions{Method: RASTERIZE}
}

type Converter interface {
	Convert(ctx context.Context, f io.Reader, opts *ConvertOptions) ([]storage.Image, error)
}

type Rotator interface {
	Rotate(ctx context.Context, image io.Reader, degrees int) (*storage.Image, error)
}

type Creator interface {
	Create(ctx context.Context, title string, records []api.Record) ([]byte, error)
}
