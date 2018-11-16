package api

import (
	"bytes"
	"context"
)

type GRPCServer struct {
	converter PdfToImageConverter
	rotator   Rotator
}

func NewGRPCServer(c PdfToImageConverter, r Rotator) *GRPCServer {
	return &GRPCServer{converter: c, rotator: r}
}

func (g *GRPCServer) ConvertPdfToImage(pdf *Pdf, sender PdfProcessor_ConvertPdfToImageServer) error {
	b := bytes.NewBuffer(pdf.Content)
	images, err := g.converter.ToImages(b)
	if err != nil {
		return err
	}

	for _, img := range images {
		err = sender.Send(img)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GRPCServer) RotateImage(ctx context.Context, rotate *Rotate) (*Image, error) {
	return g.rotator.Rotate(rotate.Content, rotate.Degree)
}
