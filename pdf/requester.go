package pdf

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"bytes"
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

func(h *HttpRequester) Do(file io.Reader) (io.ReadCloser, error) {
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

	req, err := http.NewRequest("POST", h.url, &b)
	if err != nil {
		log.Error("Error creating request")
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

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
		log.Error("PDFProcessor request error")
		return nil, err
	}
	return res.Body, nil
}
