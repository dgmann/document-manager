package pdf

import (
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"io"
)

type ImageConverter interface {
	ToImages(data io.Reader) ([]*processor.Image, error)
}

type Rotator interface {
	Rotate(img []byte, degrees float64) (*processor.Image, error)
}

type Creator interface {
	Create(document *processor.Document) (*processor.Pdf, error)
}

type PageCounter interface {
	Count(data io.Reader) (int, error)
}
