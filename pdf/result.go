package pdf

import (
	"github.com/dgmann/pdf-processor/pdfprocessor"
	"image"
	"bytes"
	log "github.com/sirupsen/logrus"
	_ "image/jpeg"
	_ "image/gif"
	_ "image/png"
)

type Result []pdfprocessor.ImageResult

func (r Result) ToImages() ([]image.Image, error) {
	images := make([]image.Image, len(r))
	for _, element := range r {
		img, s, err := image.Decode(bytes.NewReader(element.Image))
		log.Debug(s)
		if err != nil {
			log.WithField("Error", err).Panic("Error decoding image result")
			return nil, err
		}
		images[element.PageNumber] = img
	}
	return images, nil
}
