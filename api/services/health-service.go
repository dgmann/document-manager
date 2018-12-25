package services

import (
	"bytes"
	"context"
	"errors"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"net/http"
	"sync"
	"time"
)

type HealthChecker interface {
	Check() error
	CheckDb() bool
	CheckPdfProcessor() bool
}

type healthService struct {
	dbHost string
	pdfProcessorUrl string
}

var hsinstance *healthService
var hsonce sync.Once

func InitHealthService(dbHost, pdfProcessorUrl string) {
	hsonce.Do(func() {
		hsinstance = &healthService{dbHost:dbHost, pdfProcessorUrl:pdfProcessorUrl}
	})
}

func GetHealthService() *healthService {
	if hsinstance == nil {
		panic("Health Service not initialized")
	}
	return hsinstance
}

func(hs *healthService) Check() error {
	var errorString bytes.Buffer
	if !hs.CheckPdfProcessor() {
		errorString.WriteString("PdfProcessor not reachable\r\n")
	}
	if !hs.CheckDb() {
		errorString.WriteString("Database not reachable\r\n")
	}

	if errorString.Len() > 0 {
		return errors.New(errorString.String())
	}
	return nil
}

func(hs *healthService) CheckDb() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, "mongodb://"+hs.dbHost)
	if err != nil {
		return false
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return false
	}
	return true
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
