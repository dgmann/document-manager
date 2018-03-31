package pdfprocessor

import (
	"gopkg.in/gographics/imagick.v3/imagick"
	"io"
	"os/exec"
)

type ImageResult struct {
	PageNumber uint `json:"pageNumber"`
	Image []byte `json:"image"`
}



func ToImages(data io.Reader) ([]ImageResult, error) {
	fw, _ := NewFileWrapper()
	cmd := exec.Command("pdfimages", "-png" , "-tiff", "-j", "-jp2", "-jbig2", "-", fw.GetFilePath("image"))
	cmd.Stdin = data
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	images, _ := fw.GetImagesAsBuffer(".png", ".jpg")
	results := []ImageResult{}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	for i, img := range images {
		// TODO: Change resolution of png files to 50%
		results = append(results, ImageResult{PageNumber:uint(i), Image:img.Bytes()})
	}
	return results, nil
}
