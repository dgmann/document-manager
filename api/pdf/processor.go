package pdf

import (
	log "github.com/sirupsen/logrus"
	"io"
	"encoding/json"
	"bytes"
	"strconv"
	"io/ioutil"
)

type PDFProcessor struct {
	baseUrl string
}

func NewPDFProcessor(baseUrl string) *PDFProcessor {
	return &PDFProcessor{baseUrl: baseUrl}
}

func (p *PDFProcessor) Convert(f io.Reader) ([]*bytes.Buffer, error) {
	requester := NewHttpRequester(p.baseUrl + "/images/convert")
	result, err := p.Upload(requester, f)
	if err != nil {
		log.Errorf("Error fetching images: %s", err)
		return nil, err
	}
	return result.ToImages(), nil
}

func (p *PDFProcessor) Upload(requester Requester, file io.Reader) (Result, error) {
	result, err := requester.Do(file)
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

func (p *PDFProcessor) Rotate(image io.Reader, degrees int) ([]byte, error) {
	requester := NewHttpRequester(p.baseUrl + "/images/rotate/" + strconv.Itoa(degrees))
	result, err := requester.Do(image)
	if err != nil {
		log.WithField("error", err).Error("Error rotating image")
		return nil, err
	}
	defer result.Close()
	return ioutil.ReadAll(result)
}
