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

func (p *PDFProcessor) ToImages(pdf io.Reader) []image.Image {
	response := upload(p.url+"/images/extract", pdf)
	images := make([]image.Image, len(response))
	for _, element := range response {
		img, s, err := image.Decode(bytes.NewReader(element.Image))
		log.Debug(s)
		if err != nil {
			log.Panic("Error decoding image result")
		}
		images[element.PageNumber] = img
	}
	return images
}

func upload(url string, file io.Reader) []pdfprocessor.ImageResult {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	defer w.Close()

	fw, err := w.CreateFormFile("pdf", "pdf.pdf")
	if err != nil {
		log.Panic("Error creating form")
	}
	if _, err = io.Copy(fw, file); err != nil {
		log.Panic("Error copying pdf file")
	}

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		log.Panic("Error creating request")
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.WithField("error", err).Panic("Error sending request")
	}

	if res.StatusCode != http.StatusOK {
		log.WithField("status", res.Status).Panic("Request error")
	}

	var images []pdfprocessor.ImageResult
	if json.NewDecoder(res.Body).Decode(&images) != nil {
		log.WithField("error", err).Panic("Error decoding response")
	}
	return images
}
