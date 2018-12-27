package app

import (
	"context"
	"io"
)

type PdfProcessor interface {
	Converter
	Rotater
}

type Converter interface {
	Convert(ctx context.Context, f io.Reader) ([]Image, error)
}

type Rotater interface {
	Rotate(ctx context.Context, image io.Reader, degrees int) (*Image, error)
}
