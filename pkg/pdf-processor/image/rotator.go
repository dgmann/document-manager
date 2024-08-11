package image

import (
	"io"

	"github.com/dgmann/document-manager/pkg/pdf-processor/processor"
)

type Rotator interface {
	Rotate(img io.Reader, degrees float64) (*processor.Image, error)
}
