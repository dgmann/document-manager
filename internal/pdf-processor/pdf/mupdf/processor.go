package mupdf

import (
	"bytes"
	"context"
	"image/png"
	"io"

	"github.com/dgmann/document-manager/internal/pdf-processor/pdf"
	"github.com/dgmann/document-manager/pkg/processor"
	"github.com/gen2brain/go-fitz"
)

type Processor struct {
	encoder png.Encoder
}

func NewProcessor() *Processor {
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	return &Processor{encoder: encoder}
}

func (m *Processor) ToImages(ctx context.Context, data io.ReadSeeker, writer pdf.ImageSender) (int, error) {
	pdf, err := fitz.NewFromReader(data)
	if err != nil {
		return 9, err
	}
	defer pdf.Close()

	imagesSent := 0
	for i := 0; i < pdf.NumPage(); i++ {
		img, err := pdf.ImageDPI(i, 150)
		if err != nil {
			return imagesSent, err
		}
		var buf bytes.Buffer
		if err := m.encoder.Encode(&buf, img); err != nil {
			return imagesSent, err
		}
		if err := writer.Send(&processor.Image{Content: buf.Bytes(), Format: "png"}); err != nil {
			return imagesSent, err
		}
	}
	return imagesSent, nil
}

func (m *Processor) Count(data io.ReadSeeker) (int, error) {
	pdf, err := fitz.NewFromReader(data)
	if err != nil {
		return 0, err
	}
	defer pdf.Close()

	return pdf.NumPage(), nil
}
