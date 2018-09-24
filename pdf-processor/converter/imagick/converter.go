package imagick

import (
	"github.com/dgmann/document-manager/shared"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io"
	"io/ioutil"
)

type Converter struct {
}

func (c Converter) ToImages(data io.Reader) ([]*shared.Image, error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err := mw.SetResourceLimit(imagick.RESOURCE_MEMORY, 2048); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	if err := mw.SetResolution(300, 300); err != nil {
		return nil, err
	}

	pw := imagick.NewPixelWand()
	defer pw.Destroy()
	pw.SetColor("black")
	if err := mw.SetBackgroundColor(pw); err != nil {
		return nil, err
	}

	if err := mw.ReadImageBlob(b); err != nil {
		return nil, err
	}

	if err := mw.SetCompressionQuality(95); err != nil {
		return nil, err
	}
	mw.ResetIterator()

	var images []*shared.Image
	for mw.NextImage() {
		if err := mw.NormalizeImage(); err != nil {
			return nil, err
		}

		if err := mw.AutoLevelImage(); err != nil {
			return nil, err
		}

		kernel, err := imagick.NewKernelInfo("3x1:1,0,1")
		if err != nil {
			return nil, err
		}
		if err := mw.MorphologyImage(imagick.MORPHOLOGY_THICKEN, 1, kernel); err != nil {
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
		img := shared.NewImage(blob, format)
		images = append(images, img)
	}

	return images, nil
}
