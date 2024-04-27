package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/dgmann/document-manager/pdf-processor/pkg/image"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	converterFactory *pdf.ConverterFactory
	rotator          image.Rotator
	creator          pdf.Creator
	*processor.UnimplementedPdfProcessorServer
}

func NewGRPCServer(c *pdf.ConverterFactory, r image.Rotator, creator pdf.Creator) *GRPCServer {
	return &GRPCServer{converterFactory: c, rotator: r, creator: creator}
}

func (g *GRPCServer) ConvertPdfToImage(pdfFile *processor.Pdf, sender processor.PdfProcessor_ConvertPdfToImageServer) error {
	converter := g.converterFactory.Extractor()
	if pdfFile.Method == processor.Pdf_RASTERIZE {
		converter = g.converterFactory.Rasterizer()
	}

	file, err := os.CreateTemp("", "pdf-*.pdf")
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

	if _, err := io.Copy(file, bytes.NewReader(pdfFile.Content)); err != nil {
		return fmt.Errorf("error writing to temp file: %w", err)
	}
	_, _ = file.Seek(0, io.SeekStart)

	if _, err := converter.ToImages(sender.Context(), file, sender); err != nil {
		if errors.Is(err, pdf.ErrorExtraction) {
			return status.Error(400, err.Error())
		}
		return err
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
