package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/datastore"
	"github.com/sirupsen/logrus"
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

func NewHttpUploader(url string, timeout time.Duration) *HttpUploader {
	return &HttpUploader{url, &http.Client{
		Timeout: timeout,
	}}
}

type NewRecord struct {
	File       io.Reader
	ReceivedAt time.Time
	Sender     string
	Date       *time.Time
	PatientId  *string
	Status     *datastore.Status
	Comment    *string
	Category   *string
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
	if res.StatusCode == http.StatusCreated {
		return nil
	}
	resBody, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == http.StatusConflict {
		logrus.Infof("Category %s already exists. Skipping: %s", category.Id, string(resBody))
		return nil
	}
	return errors.New(fmt.Sprintf("bad response. Status: %v, Message: %s", res.StatusCode, string(resBody)))
}

func (u *HttpUploader) CreateRecord(create *NewRecord) error {
	params := createParamMap(create)
	req, err := newfileUploadRequest(u.url+"/records", params, create.File)
	if err != nil {
		return err
	}
	resp, err := u.client.Do(req)
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
	if !create.Date.IsZero() {
		params["date"] = create.Date.Format(time.RFC3339)
	}
	if !create.ReceivedAt.IsZero() {
		params["receivedAt"] = create.ReceivedAt.Format(time.RFC3339)
	}
	if create.PatientId != nil {
		params["patientId"] = *create.PatientId
	}
	if create.Status != nil {
		params["status"] = string(*create.Status)
	}
	if create.Comment != nil {
		params["comment"] = *create.Comment
	}
	if create.Category != nil {
		params["category"] = *create.Category
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
