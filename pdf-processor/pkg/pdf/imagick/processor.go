package imagick

import (
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io"
	"io/ioutil"
	"strings"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (c *Processor) Rotate(img []byte, degrees float64) (*processor.Image, error) {
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
	return &processor.Image{Content: blob, Format: format}, nil
}

func (c *Processor) ToImages(data io.Reader) ([]*processor.Image, error) {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.SetResolution(300, 300); err != nil {
		return nil, err
	}

	if err := mw.ReadImageBlob(b); err != nil {
		return nil, err
	}

	pw := imagick.NewPixelWand()
	defer pw.Destroy()
	pw.SetColor("white")
	if err := mw.SetBackgroundColor(pw); err != nil {
		return nil, err
	}

	if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_DEACTIVATE); err != nil {
		return nil, err
	}

	if err := mw.SetCompressionQuality(95); err != nil {
		return nil, err
	}
	mw.ResetIterator()

	var images []*processor.Image
	for mw.NextImage() {
		if err := mw.NormalizeImage(); err != nil {
			return nil, err
		}

		if err := mw.AutoLevelImage(); err != nil {
			return nil, err
		}

		factor := float64(mw.GetImageWidth()) / 720.0
		if err := mw.ScaleImage(uint(float64(mw.GetImageWidth())/factor), uint(float64(mw.GetImageHeight())/factor)); err != nil {
			return nil, err
		}

		format := "png"
		if err := mw.SetFormat(format); err != nil {
			return nil, err
		}

		if err := mw.QuantizeImage(256, imagick.COLORSPACE_RGB, 1, imagick.DITHER_METHOD_UNDEFINED, false); err != nil {
			return nil, err
		}

		if err := mw.SetImageDepth(8); err != nil {
			return nil, err
		}

		blob := mw.GetImageBlob()

		img := &processor.Image{Content: blob, Format: format}
		images = append(images, img)
	}

	return images, nil
}
