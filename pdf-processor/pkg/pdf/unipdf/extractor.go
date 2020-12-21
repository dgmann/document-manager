package unipdf

import (
	"image"
	"image/png"
	"io"
	"sync"

	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pool"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	gim "github.com/ozankasikci/go-image-merge"
	"github.com/unidoc/unipdf/v3/extractor"
	unipdf "github.com/unidoc/unipdf/v3/model"
)

type syncpool struct{ sync.Pool }

type Extractor struct {
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (e *Extractor) ToImages(seeker io.ReadSeeker, writer pdf.ImageSender) (int, error) {
	pdfReader, err := unipdf.NewPdfReader(seeker)
	if err != nil {
		return 0, err
	}
	imagesSent := 0
	for _, page := range pdfReader.PageList {
		img, err := extractImage(page)
		if err != nil {
			return imagesSent, err
		}

		if err := writer.Send(&processor.Image{Content: img, Format: "png"}); err != nil {
			return 0, err
		}
		imagesSent++
	}
	return imagesSent, nil
}

func extractImage(page *unipdf.PdfPage) ([]byte, error) {
	extract, err := extractor.New(page)
	if err != nil {
		return nil, err
	}
	images, err := extract.ExtractPageImages(nil)
	if err != nil {
		return nil, err
	}
	if len(images.Images) == 1 {
		goImg, err := images.Images[0].Image.ToGoImage()
		if err != nil {
			return nil, err
		}
		return encode(goImg)
	}
	goImages := make([]image.Image, len(images.Images))
	for i, img := range images.Images {
		goImg, err := img.Image.ToGoImage()
		if err != nil {
			return nil, err
		}
		goImages[i] = goImg
	}
	return concatenateImages(goImages)
}

func concatenateImages(imgs []image.Image) ([]byte, error) {
	grids := make([]*gim.Grid, 0, len(imgs))
	for _, img := range imgs {
		img2 := img
		grids = append(grids, &gim.Grid{Image: &img2})
	}

	img, err := gim.New(grids, 1, len(imgs)).Merge()
	if err != nil {
		return nil, err
	}

	return encode(img)
}

func encode(img image.Image) ([]byte, error) {
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
