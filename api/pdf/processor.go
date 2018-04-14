package pdf

import (
	log "github.com/sirupsen/logrus"
	"io"
	"encoding/json"
	"strconv"
	"github.com/dgmann/document-manager/shared"
)

type PDFProcessor struct {
	baseUrl string
}

func NewPDFProcessor(baseUrl string) *PDFProcessor {
	return &PDFProcessor{baseUrl: baseUrl}
}

func (p *PDFProcessor) Convert(f io.Reader) ([]*shared.Image, error) {
	requester := NewHttpRequester(p.baseUrl + "/images/convert")
	result, err := p.Upload(requester, f)
	if err != nil {
		log.Errorf("Error fetching images: %s", err)
		return nil, err
	}
	return result, nil
}

func (p *PDFProcessor) Upload(requester Requester, file io.Reader) ([]*shared.Image, error) {
	result, err := requester.Do(file)
	if err != nil {
		log.WithField("error", err).Error("Error transforming pdf to images")
		return nil, err
	}
	defer result.Close()
	var images []*shared.Image
	if json.NewDecoder(result).Decode(&images) != nil {
		log.WithField("error", err).Error("Error decoding response")
		return nil, err
	}
	return images, nil
}

func (p *PDFProcessor) Rotate(image io.Reader, degrees int) (*shared.Image, error) {
	requester := NewHttpRequester(p.baseUrl + "/images/rotate/" + strconv.Itoa(degrees))
	result, err := requester.Do(image)
	if err != nil {
		log.WithField("error", err).Error("Error rotating image")
		return nil, err
	}
	defer result.Close()
	var img shared.Image
	if json.NewDecoder(result).Decode(&img) != nil {
		log.WithField("error", err).Error("Error decoding response")
		return nil, err
	}
	return &img, nil
}
