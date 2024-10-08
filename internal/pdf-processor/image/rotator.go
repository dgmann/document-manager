package image

import (
	"io"

	"github.com/dgmann/document-manager/pkg/processor"
)

type Rotator interface {
	Rotate(img io.Reader, degrees float64) (*processor.Image, error)
}
