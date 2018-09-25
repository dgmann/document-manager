package imagick

import (
	"github.com/dgmann/document-manager/pdf-processor/api"
	"gopkg.in/gographics/imagick.v3/imagick"
	"strings"
)

func (c *Processor) Rotate(img []byte, degrees float64) (*api.Image, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImageBlob(img)
	if err != nil {
		return nil, err
	}
	format := strings.ToLower(mw.GetImageFormat())

	pw := imagick.NewPixelWand()
	defer pw.Destroy()

	pw.SetColor("black")
	err = mw.RotateImage(pw, degrees)
	if err != nil {
		return nil, err
	}
	mw.ResetIterator()
	blob := mw.GetImageBlob()
	return &api.Image{Content: blob, Format: format}, nil
}
