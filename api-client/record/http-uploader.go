package record

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type HttpUploader struct {
	url string
}

func NewHttpUploader(url string) *HttpUploader {
	return &HttpUploader{url}
}

func (u *HttpUploader) Upload(create *NewRecord) error {
	params := createParamMap(create)
	req, err := newfileUploadRequest(u.url+"/records", params, create.File)
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

func createParamMap(create *NewRecord) map[string]string {
	params := map[string]string{
		"sender": create.Sender,
	}
	if create.Date != nil {
		params["date"] = create.Date.String()
	}
	if create.ReceivedAt != nil {
		params["receivedAt"] = create.ReceivedAt.String()
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
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
