package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient struct {
	client     *httpClient
	Records    *RecordClient
	Categories *CategoryClient
}

func NewHTTPClient(uri string, timeout time.Duration) (*HTTPClient, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	c := &httpClient{
		Url: u,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
	records := &RecordClient{
		httpClient: c,
	}
	categories := &CategoryClient{
		httpClient: c,
	}
	return &HTTPClient{client: c, Records: records, Categories: categories}, nil
}

type Category struct {
	Id   string `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name,omitempty" json:"name"`
}

var (
	ErrAlreadyExists = errors.New("resource already exists")
	ErrNotFound      = errors.New("resource not found")
	ErrGeneric       = errors.New("error handling response")
)

type httpClient struct {
	Url    *url.URL
	Client *http.Client
}

func (c *httpClient) PostJson(endpoint string, payload interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, fmt.Errorf("error encoding body to json: %w", err)
	}

	p := c.Url.JoinPath(endpoint)
	return c.Client.Post(p.String(), "application/json", body)
}

func (c *httpClient) PutJson(endpoint string, payload interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, fmt.Errorf("error encoding body to json: %w", err)
	}

	p := c.Url.JoinPath(endpoint)
	req, err := http.NewRequest(http.MethodPut, p.String(), body)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	return resp, nil
}

func (c *httpClient) GetJson(endpoint string) (*http.Response, error) {
	p := c.Url.JoinPath(endpoint)
	return c.Client.Get(p.String())
}

func HandleResponse[T any](res *http.Response) (T, error) {
	var result T
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		if err := ParseJsonBody(res.Body, &result); err != nil {
			return result, fmt.Errorf("error parsing json body: %w", err)
		}
		return result, nil
	}
	var resBody map[string]interface{}
	if err := ParseJsonBody(res.Body, &resBody); err != nil {
		return result, fmt.Errorf("error parsing json body: %w", err)
	}
	if res.StatusCode == http.StatusNotFound {
		return result, fmt.Errorf("%w: %+v", ErrNotFound, resBody)
	}
	if res.StatusCode == http.StatusConflict {
		return result, fmt.Errorf("%w: %+v", ErrAlreadyExists, resBody)
	}

	return result, fmt.Errorf("%w: status %d, %+v", ErrGeneric, res.StatusCode, resBody)
}

func ParseJsonBody(res io.ReadCloser, out interface{}) (err error) {
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if err == nil {
			err = closeErr
		}
	}(res)
	return json.NewDecoder(res).Decode(out)
}

func ToPointer[T any](in T, err error) (*T, error) {
	return &in, err
}
