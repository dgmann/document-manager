package mupdf

import (
	"bytes"
	"image/png"
	"io"

	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/gen2brain/go-fitz"
)

type Processor struct {
	encoder png.Encoder
}

func NewProcessor() *Processor {
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	return &Processor{encoder: encoder}
}

func (m *Processor) ToImages(data io.Reader) ([]*processor.Image, error) {
	pdf, err := fitz.NewFromReader(data)
	if err != nil {
		return nil, err
	}
	defer pdf.Close()

	images := make([]*processor.Image, pdf.NumPage())
	for i := 0; i < pdf.NumPage(); i++ {
		img, err := pdf.ImageDPI(i, 150)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		if err := m.encoder.Encode(&buf, img); err != nil {
			return nil, err
		}
		images[i] = &processor.Image{Content: buf.Bytes(), Format: "png"}
	}
	return images, nil
}

func (m *Processor) Count(data io.Reader) (int, error) {
	pdf, err := fitz.NewFromReader(data)
	if err != nil {
		return 0, err
	}
	defer pdf.Close()

	return pdf.NumPage(), nil
}
