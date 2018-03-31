package pdfprocessor

import (
	"gopkg.in/gographics/imagick.v3/imagick"
	"github.com/dgmann/document-manager/shared"
	"io"
	"os/exec"
)

func ToImages(data io.Reader) ([]shared.ImageResult, error) {
	fw, _ := NewFileWrapper()
	cmd := exec.Command("pdfimages", "-png" , "-tiff", "-j", "-jp2", "-jbig2", "-", fw.GetFilePath("image"))
	cmd.Stdin = data
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	images, _ := fw.GetImagesAsBuffer(".png", ".jpg")
	var results []shared.ImageResult

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	for i, img := range images {
		// TODO: Change resolution of png files to 50%
		results = append(results, shared.ImageResult{PageNumber:uint(i), Image:img.Bytes()})
	}
	return results, nil
}
