package pdfprocessor

import (
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
	return images, nil
}
