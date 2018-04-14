package pdfprocessor

import (
	"gopkg.in/gographics/imagick.v3/imagick"
	"github.com/dgmann/document-manager/shared"
	"io"
	"os/exec"
	log "github.com/sirupsen/logrus"
)

func ToImages(data io.Reader) ([]*shared.Image, error) {
	fw, err := NewFileWrapper()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer fw.Destroy()
	cmd := exec.Command("pdfimages", "-png" , "-tiff", "-j", "-jp2", "-jbig2", "-", fw.GetFilePath("image"))
	cmd.Stdin = data
	err = cmd.Run()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	images, _ := fw.GetImagesAsBuffer(".png", ".jpg")
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	for _, img := range images {
		mw.ReadImageBlob(img.Image)

		width := mw.GetImageWidth()
		height := mw.GetImageHeight()
		if height > 2000 {
			err := mw.AdaptiveResizeImage(width/2, height/2)
			if err != nil {
				return nil, err
			}
			mw.ResetIterator()
			img.Image = mw.GetImageBlob()
		}
	}
	return images, nil
}
