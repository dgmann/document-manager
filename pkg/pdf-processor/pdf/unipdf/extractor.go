package unipdf

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"sync"

	"github.com/dgmann/document-manager/pkg/pdf-processor/pdf"
	"github.com/dgmann/document-manager/pkg/pdf-processor/pool"
	"github.com/dgmann/document-manager/pkg/pdf-processor/processor"
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

func (e *Extractor) ToImages(ctx context.Context, seeker io.ReadSeeker, writer pdf.ImageSender) (int, error) {
	pdfReader, err := unipdf.NewPdfReader(seeker)
	if err != nil {
		return 0, err
	}
	imagesSent := 0
	for _, page := range pdfReader.PageList {
		if err := convertPage(writer, page); err != nil {
			return imagesSent, err
		}
		imagesSent++
	}
	return imagesSent, nil
}

func convertPage(writer pdf.ImageSender, page *unipdf.PdfPage) error {
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)

	if err := extractImage(buf, page); err != nil {
		return err
	}

	if err := writer.Send(&processor.Image{Content: buf.Bytes(), Format: "png"}); err != nil {
		return err
	}
	return nil
}

func extractImage(buf io.Writer, page *unipdf.PdfPage) error {
	extract, err := extractor.New(page)
	if err != nil {
		return err
	}
	images, err := extract.ExtractPageImages(nil)
	if err != nil {
		return err
	}

	if box, err := page.GetMediaBox(); err == nil {
		for _, img := range images.Images {
			if box.Width()/img.X < 3 {
				return fmt.Errorf("starting position of found image is not in the first third of the page. page width: %d, image X: %d", int64(box.Width()), int64(img.X))
			}
		}
	}

	if len(images.Images) == 1 {
		goImg, err := images.Images[0].Image.ToGoImage()
		if err != nil {
			return err
		}
		return encode(buf, goImg)
	}
	goImages := make([]image.Image, len(images.Images))
	for i, img := range images.Images {
		goImg, err := img.Image.ToGoImage()
		if err != nil {
			return err
		}
		goImages[i] = goImg
	}
	return concatenateImages(buf, goImages)
}

func concatenateImages(buf io.Writer, imgs []image.Image) error {
	grids := make([]*gim.Grid, 0, len(imgs))
	for _, img := range imgs {
		img2 := img
		grids = append(grids, &gim.Grid{Image: &img2})
	}

	img, err := gim.New(grids, 1, len(imgs)).Merge()
	if err != nil {
		return err
	}

	return encode(buf, img)
}

func encode(buf io.Writer, img image.Image) error {
	if err := png.Encode(buf, img); err != nil {
		return err
	}
	return nil
}
