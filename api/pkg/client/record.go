package client

import (
	"bytes"
	"fmt"
	"github.com/dgmann/document-manager/api/pkg/api"
	"io"
	"mime/multipart"
	"time"
)

type NewRecord struct {
	File       io.Reader
	ReceivedAt time.Time
	Sender     string
	Date       *time.Time
	PatientId  *string
	Status     *api.Status
	Comment    *string
	Category   *string
}

type RecordClient struct {
	*httpClient
}

func (c *RecordClient) Get(id string) (*api.Record, error) {
	res, err := c.GetJson("records/" + id)
	if err != nil {
		return nil, err
	}
	return ToPointer(HandleResponse[api.Record](res))
}

func (c *RecordClient) Create(create *NewRecord) (*api.Record, error) {
	params := createParamMap(create)
	body, contentType, err := newfileUploadRequest(params, create.File)
	if err != nil {
		return nil, err
	}

	reqUrl := c.Url.JoinPath("records").String()
	res, err := c.Client.Post(reqUrl, contentType, body)
	if err != nil {
		return nil, err
	}

	return ToPointer(HandleResponse[api.Record](res))
}

func (c *RecordClient) Update(record *api.Record) (*api.Record, error) {
	res, err := c.PutJson(fmt.Sprintf("records/%s", record.Id), record)
	if err != nil {
		return nil, err
	}

	return ToPointer(HandleResponse[api.Record](res))
}

func (c *RecordClient) UpdatePages(recordId string, updatedPages []api.PageUpdate) (*api.Record, error) {
	res, err := c.PutJson(fmt.Sprintf("records/%s/pages", recordId), updatedPages)
	if err != nil {
		return nil, err
	}

	return ToPointer(HandleResponse[api.Record](res))
}

func createParamMap(create *NewRecord) map[string]string {
	params := map[string]string{
		"sender": create.Sender,
	}
	if create.Date != nil && !create.Date.IsZero() {
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

func newfileUploadRequest(params map[string]string, file io.Reader) (b io.Reader, contentType string, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer func(writer *multipart.Writer) {
		closeErr := writer.Close()
		if err == nil {
			err = closeErr
		}
	}(writer)

	part, err := writer.CreateFormFile("pdf", "file.pdf")
	if err != nil {
		return body, writer.FormDataContentType(), fmt.Errorf("error creating form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return body, writer.FormDataContentType(), fmt.Errorf("error copying file to form field: %w", err)
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	return body, writer.FormDataContentType(), nil
}
