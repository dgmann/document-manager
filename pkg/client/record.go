package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/dgmann/document-manager/pkg/api"
)

type RecordClient struct {
	*httpClient
}

func (c *RecordClient) List() ([]api.Record, error) {
	res, err := c.GetJson("records")
	if err != nil {
		return nil, err
	}
	return HandleResponse[[]api.Record](res)
}

func (c *RecordClient) Get(id string) (*api.Record, error) {
	res, err := c.GetJson("records/" + id)
	if err != nil {
		return nil, err
	}
	return ToPointer(HandleResponse[api.Record](res))
}

func (c *RecordClient) Create(create *api.NewRecord) (*api.Record, error) {
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

func (c *RecordClient) Download(recordId string) (io.ReadCloser, error) {
	p := c.Url.JoinPath("export")
	query := p.Query()
	query.Add("id", recordId)
	p.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodGet, p.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		var resBody map[string]interface{}
		if err := ParseJsonBody(res.Body, &resBody); err != nil {
			return nil, fmt.Errorf("error parsing json body: %w", err)
		}
		return nil, fmt.Errorf("error downloading file, status: %d, body: %+v", res.StatusCode, resBody)
	}
	return res.Body, nil
}

func createParamMap(create *api.NewRecord) map[string]string {
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
