package app

import (
	"context"
	"io"
)

type PdfProcessor interface {
	Converter
	Rotater
	PdfCreator
}

type Converter interface {
	Convert(ctx context.Context, f io.Reader) ([]Image, error)
}

type Rotater interface {
	Rotate(ctx context.Context, image io.Reader, degrees int) (*Image, error)
}

type PdfCreator interface {
	CreatePdf(ctx context.Context, title string, records []Record) ([]byte, error)
}
