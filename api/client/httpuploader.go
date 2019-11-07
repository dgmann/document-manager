package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/datastore"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type HttpUploader struct {
	url    string
	client *http.Client
}

func NewHttpUploader(url string) *HttpUploader {
	return &HttpUploader{url, &http.Client{
		Timeout: time.Second * 10,
	}}
}

type NewRecord struct {
	File       io.Reader
	ReceivedAt time.Time
	Sender     string
}

func (u *HttpUploader) CreateCategory(category datastore.Category) error {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(category); err != nil {
		return err
	}

	res, err := http.Post(u.url+"/categories", "application/json", body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (u *HttpUploader) CreateRecord(create *NewRecord) error {
	params := createParamMap(create)
	req, err := newfileUploadRequest(u.url+"/records", params, create.File)
	if err != nil {
		return err
	}
	client := u.client
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("bad response. Status: %v, Message: %s", resp.StatusCode, string(body)))
	}
	return nil
}

func createParamMap(create *NewRecord) map[string]string {
	params := map[string]string{
		"sender": create.Sender,
	}
	if !create.ReceivedAt.IsZero() {
		params["receivedAt"] = create.ReceivedAt.Format(time.RFC3339)
	}
	return params
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
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
