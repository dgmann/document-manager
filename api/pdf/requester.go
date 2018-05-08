package pdf

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"io"
	"errors"
	"fmt"
	"io/ioutil"
)

type Requester interface {
	Do(b io.Reader) (io.ReadCloser, error)
}

type HttpRequester struct {
	url string
	client *http.Client
}

func NewHttpRequester(url string) *HttpRequester {
	return &HttpRequester{url: url, client: &http.Client{}}
}

func(h *HttpRequester) Do(file io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", h.url, file)
	if err != nil {
		log.Error("Error creating request")
		return nil, err
	}

	res, err := h.client.Do(req)
	if err != nil {
		log.WithField("error", err).Error("Error sending request")
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		bodyString := string(bodyBytes)
		log.WithFields(log.Fields{
			"status": res.Status,
			"error": bodyString,
		}).Error("Request error")
		return nil, errors.New(fmt.Sprintf("processor responded with error code %d: %s", res.StatusCode, bodyString))
	}

	if err != nil {
		log.Error("Processor request error")
		return nil, err
	}
	return res.Body, nil
}
