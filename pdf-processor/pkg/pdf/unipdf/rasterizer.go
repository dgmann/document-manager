package unipdf

import (
	"io"

	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	pdf "github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
)

type Rasterizer struct {
}

func NewRasterizer() *Rasterizer {
	return &Rasterizer{}
}

func (e *Rasterizer) ToImages(data io.ReadSeeker) ([]*processor.Image, error) {
	pdfReader, err := pdf.NewPdfReader(data)
	if err != nil {
		return nil, err
	}
	renderer := render.NewImageDevice()
	var images []*processor.Image
	for _, page := range pdfReader.PageList {
		res, err := renderer.Render(page)
		if err != nil {
			return nil, err
		}
		img, err := encode(res)
		if err != nil {
			return nil, err
		}
		images = append(images, &processor.Image{Content: img, Format: "png"})
	}
	return images, nil
}
