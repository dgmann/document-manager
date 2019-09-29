package main

import (
	"bytes"
	"context"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"io"
)

type PdfToImageConverter interface {
	ToImages(data io.Reader) ([]*processor.Image, error)
}

type Rotator interface {
	Rotate(img []byte, degrees float64) (*processor.Image, error)
}

type Creator interface {
	Create(document *processor.Document) (*processor.Pdf, error)
}

type GRPCServer struct {
	converter PdfToImageConverter
	rotator   Rotator
	creator   Creator
}

func NewGRPCServer(c PdfToImageConverter, r Rotator, creator Creator) *GRPCServer {
	return &GRPCServer{converter: c, rotator: r, creator: creator}
}

func (g *GRPCServer) ConvertPdfToImage(pdf *processor.Pdf, sender processor.PdfProcessor_ConvertPdfToImageServer) error {
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

func (g *GRPCServer) RotateImage(ctx context.Context, rotate *processor.Rotate) (*processor.Image, error) {
	return g.rotator.Rotate(rotate.Content, rotate.Degree)
}

func (g *GRPCServer) CreatePdf(ctx context.Context, document *processor.Document) (*processor.Pdf, error) {
	return g.creator.Create(document)
}
