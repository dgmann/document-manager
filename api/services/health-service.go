package services

import (
	"sync"
	"net/http"
	"time"
	"errors"
)

type healthService struct {
	pdfProcessorUrl string
}

var hsinstance *healthService
var hsonce sync.Once

func InitHealthService(pdfProcessorUrl string) {
	hsonce.Do(func() {
		hsinstance = &healthService{pdfProcessorUrl:pdfProcessorUrl}
	})
}

func GetHealthService() *healthService {
	if hsinstance == nil {
		panic("Health Service not initialized")
	}
	return hsinstance
}

func(hs *healthService) Check() error {
	if hs.CheckPdfProcessor() {
		return nil
	}
	return errors.New("PdfProcessor not reachable")
}

func(hs *healthService) CheckPdfProcessor() bool {
	return httpCheck(hs.pdfProcessorUrl)
}

func httpCheck(url string) bool {
	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	res, err := client.Do(req)
	if err != nil {
		return false
	}
	return res.StatusCode == 200
}
