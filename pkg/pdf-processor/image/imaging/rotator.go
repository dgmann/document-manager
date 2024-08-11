package imaging

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/dgmann/document-manager/pkg/pdf-processor/processor"
	"github.com/disintegration/imaging"
)

type Rotator struct{}

func NewRotator() *Rotator {
	return &Rotator{}
}

func (r *Rotator) Rotate(data io.Reader, degrees float64) (*processor.Image, error) {
	img, format, err := image.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}
	rotated := imaging.Rotate(img, degrees, color.Transparent)
	var buf bytes.Buffer
	imgagingFormat, err := imaging.FormatFromExtension(format)
	if err != nil {
		return nil, fmt.Errorf("error getting imaging format: %w", err)
	}
	if err := imaging.Encode(&buf, rotated, imgagingFormat); err != nil {
		return nil, fmt.Errorf("error encoding image: %w", err)
	}
	return &processor.Image{Content: buf.Bytes(), Format: format}, nil
}
