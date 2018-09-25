package api

import "io"

type PdfToImageConverter interface {
	ToImages(data io.Reader) ([]*Image, error)
}
