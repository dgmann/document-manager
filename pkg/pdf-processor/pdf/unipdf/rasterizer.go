package unipdf

import (
	"context"
	"errors"
	"io"

	"github.com/dgmann/document-manager/pkg/pdf-processor/pdf"
	"github.com/dgmann/document-manager/pkg/pdf-processor/pool"
	"github.com/dgmann/document-manager/pkg/pdf-processor/processor"
	unipdf "github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
)

type Rasterizer struct {
}

func NewRasterizer() *Rasterizer {
	return &Rasterizer{}
}

func (e *Rasterizer) ToImages(ctx context.Context, data io.ReadSeeker, writer pdf.ImageSender) (int, error) {
	pdfReader, err := unipdf.NewPdfReader(data)
	if err != nil {
		return 0, err
	}
	renderer := render.NewImageDevice()
	imagesSent := 0
	for _, page := range pdfReader.PageList {
		if errors.Is(ctx.Err(), context.Canceled) {
			return imagesSent, ctx.Err()
		}
		if err := rasterizePage(writer, renderer, page); err != nil {
			return imagesSent, err
		}
	}
	return imagesSent, nil
}

func rasterizePage(writer pdf.ImageSender, renderer *render.ImageDevice, page *unipdf.PdfPage) error {
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)

	if err := rasterize(renderer, buf, page); err != nil {
		return err
	}

	if err := writer.Send(&processor.Image{Content: buf.Bytes(), Format: "png"}); err != nil {
		return err
	}
	return nil
}

func rasterize(renderer *render.ImageDevice, buf io.Writer, page *unipdf.PdfPage) error {
	res, err := renderer.Render(page)
	if err != nil {
		return err
	}
	if encode(buf, res) != nil {
		return err
	}
	return nil
}
