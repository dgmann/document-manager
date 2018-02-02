package services

import (
	"bytes"
	"encoding/json"
	"github.com/dgmann/pdf-processor/pdfprocessor"
	log "github.com/sirupsen/logrus"
	"image"
	"io"
	"mime/multipart"
	"net/http"
)

type PDFProcessor struct {
	url string
}

func NewPDFProcessor(url string) *PDFProcessor {
	return &PDFProcessor{url: url}
}

func (p *PDFProcessor) ToImages(pdf io.Reader) ([]image.Image, error) {
	response, err := upload(p.url+"/images/extract", pdf)
	if err != nil {
		return nil, err
	}
	images := make([]image.Image, len(response))
	for _, element := range response {
		img, s, err := image.Decode(bytes.NewReader(element.Image))
		log.Debug(s)
		if err != nil {
			log.Panic("Error decoding image result")
		}
		images[element.PageNumber] = img
	}
	return images, nil
}

func upload(url string, file io.Reader) ([]pdfprocessor.ImageResult, error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("pdf", "pdf.pdf")
	if err != nil {
		log.Error("Error creating form")
		return nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		log.Error("Error copying pdf file")
		return nil, err
	}
	if w.Close() != nil {
		log.Error("Error closing multipart writer")
		return nil, err
	}

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		log.Error("Error creating request")
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.WithField("error", err).Error("Error sending request")
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.WithField("status", res.Status).Error("Request error")
		return nil, err
	}

	var images []pdfprocessor.ImageResult
	if json.NewDecoder(res.Body).Decode(&images) != nil {
		log.WithField("error", err).Error("Error decoding response")
		return nil, err
	}
	return images, nil
}
