package image

import (
	"gopkg.in/gographics/imagick.v3/imagick"
	"strings"
)

func Rotate(img []byte, degrees float64) ([]byte, string, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	err := mw.ReadImageBlob(img)
	if err != nil {
		return nil, "", err
	}
	format := strings.ToLower(mw.GetImageFormat())

	pw := imagick.NewPixelWand()
	defer pw.Destroy()
	pw.SetColor("black")

	err = mw.RotateImage(pw, degrees)
	if err != nil {
		return nil, format, err
	}
	mw.ResetIterator()
	return mw.GetImageBlob(), format, nil
}
