package upload

import (
	"github.com/dgmann/document-manager/directory-watcher/models"
	"net/http"
	"bytes"
	"mime/multipart"
	"io"
	"errors"
	"fmt"
	"os"
	"io/ioutil"
)

type HttpUploader struct {
	url string
}

func NewHttpUploader(url string) *HttpUploader {
	return &HttpUploader{url}
}

func (u *HttpUploader) Upload(create *models.RecordCreate) error {
	params := map[string]string{
		"sender":     create.Sender,
		"receivedAt": create.ReceivedAt.String(),
	}
	pdf, err := os.Open(create.PdfPath)
	if err != nil {
		return err
	}
	req, err := newfileUploadRequest(u.url+"/records", params, pdf)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("bad response. Status: %v, Message: %s", resp.StatusCode, string(body)))
	}
	return nil
}

func newfileUploadRequest(uri string, params map[string]string, file io.Reader) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("pdf", "file.pdf")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
