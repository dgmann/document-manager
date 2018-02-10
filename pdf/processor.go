package pdf

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"encoding/json"
	"bytes"
)

type PDFProcessor struct {
	requester Requester
}

func NewPDFProcessor(url string) *PDFProcessor {
	return &PDFProcessor{requester: &HttpRequester{url: url + "/images/extract", client: &http.Client{}}}
}

func (p *PDFProcessor) Convert(f io.Reader) ([]*bytes.Buffer, error) {
	result, err := p.Upload(f)
	if err != nil {
		log.Errorf("Error fetching images: %s", err)
		return nil, err
	}
	return result.ToImages(), nil
}

func (p * PDFProcessor) Upload(file io.Reader) (Result, error) {
	result, err := p.requester.Do(file)
	if err != nil {
		log.WithField("error", err).Error("Error transforming pdf to images")
		return nil, err
	}
	defer result.Close()
	var images Result
	if json.NewDecoder(result).Decode(&images) != nil {
		log.WithField("error", err).Error("Error decoding response")
		return nil, err
	}
	return images, nil
}
