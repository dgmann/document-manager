package converter

import (
	"github.com/dgmann/document-manager/shared"
	"io"
)

type PdfToImageConverter interface {
	ToImages(data io.Reader) ([]*shared.Image, error)
}
