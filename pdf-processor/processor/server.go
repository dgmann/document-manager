package processor

import (
	"bytes"
	"context"
	"github.com/dgmann/document-manager/pdf-processor/converter"
	"github.com/dgmann/document-manager/pdf-processor/image"
)

type GRPCServer struct {
	converter converter.PdfToImageConverter
}

func NewGRPCServer(c converter.PdfToImageConverter) *GRPCServer {
	return &GRPCServer{converter: c}
}

func (g *GRPCServer) ConvertPdfToImage(pdf *Pdf, sender PdfProcessor_ConvertPdfToImageServer) error {
	b := bytes.NewBuffer(pdf.Content)
	images, err := g.converter.ToImages(b)
	if err != nil {
		return err
	}

	for _, img := range images {
		sender.Send(&Image{
			Content: img.Image,
			Format:  img.Format,
		})
	}

	return nil
}

func (g *GRPCServer) RotateImage(ctx context.Context, rotate *Rotate) (*Image, error) {
	rotated, err := image.Rotate(rotate.Content, rotate.Degree)
	if err != nil {
		return nil, err
	}
	return &Image{
		Content: rotated.Image,
		Format:  rotated.Format,
	}, nil
}
