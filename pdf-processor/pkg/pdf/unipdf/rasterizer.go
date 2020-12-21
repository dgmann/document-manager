package unipdf

import (
	"io"

	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	unipdf "github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
)

type Rasterizer struct {
}

func NewRasterizer() *Rasterizer {
	return &Rasterizer{}
}

func (e *Rasterizer) ToImages(data io.ReadSeeker, writer pdf.ImageSender) (int, error) {
	pdfReader, err := unipdf.NewPdfReader(data)
	if err != nil {
		return 0, err
	}
	renderer := render.NewImageDevice()
	imagesSent := 0
	for _, page := range pdfReader.PageList {
		res, err := renderer.Render(page)
		if err != nil {
			return imagesSent, err
		}
		img, err := encode(res)
		if err != nil {
			return imagesSent, err
		}
		if err := writer.Send(&processor.Image{Content: img, Format: "png"}); err != nil {
			return imagesSent, err
		}
	}
	return imagesSent, nil
}
