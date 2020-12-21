package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/dgmann/document-manager/pdf-processor/pkg/image"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/sirupsen/logrus"
)

type GRPCServer struct {
	converter pdf.ImageConverter
	rotator   image.Rotator
	creator   pdf.Creator
	*processor.UnimplementedPdfProcessorServer
}

func NewGRPCServer(c pdf.ImageConverter, r image.Rotator, creator pdf.Creator) *GRPCServer {
	return &GRPCServer{converter: c, rotator: r, creator: creator}
}

func (g *GRPCServer) ConvertPdfToImage(pdf *processor.Pdf, sender processor.PdfProcessor_ConvertPdfToImageServer) error {
	file, err := ioutil.TempFile("", "pdf-*.pdf")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logrus.Warn(err)
		}
		if err := os.Remove(file.Name()); err != nil {
			logrus.Warn(err)
		}
	}()

	if _, err := io.Copy(file, bytes.NewReader(pdf.Content)); err != nil {
		return fmt.Errorf("error writing to temp file: %w", err)
	}
	_, _ = file.Seek(0, io.SeekStart)

	images, err := g.converter.ToImages(file)
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
	buf := bytes.NewBuffer(rotate.Content)
	return g.rotator.Rotate(buf, rotate.Degree)
}

func (g *GRPCServer) CreatePdf(ctx context.Context, document *processor.Document) (*processor.Pdf, error) {
	return g.creator.Create(document)
}
