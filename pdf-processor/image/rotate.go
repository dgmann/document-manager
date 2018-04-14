package image

import (
	"gopkg.in/gographics/imagick.v3/imagick"
	"strings"
	"github.com/dgmann/document-manager/shared"
)

func Rotate(img []byte, degrees float64) (*shared.Image, error) {
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
	return shared.NewImage(blob, format), nil
}
