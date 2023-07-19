package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient struct {
	client  *httpClient
	Records *RecordClient
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
	return &HTTPClient{client: c, Records: records}, nil
}

type Category struct {
	Id   string `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name,omitempty" json:"name"`
}

var ErrAlreadyExists = errors.New("resource already exists")

type httpClient struct {
	Url    *url.URL
	Client *http.Client
}

func (c *httpClient) PostJson(endpoint string, payload interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	p := c.Url.JoinPath(endpoint)
	return c.Client.Post(p.String(), "application/json", body)
}

func (c *httpClient) PutJson(endpoint string, payload interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	p := c.Url.JoinPath(endpoint)
	req, err := http.NewRequest(http.MethodPut, p.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.Client.Do(req)
}

func (c *httpClient) GetJson(endpoint string, out interface{}) error {
	p := c.Url.JoinPath(endpoint)
	res, err := c.Client.Get(p.String())
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}(res.Body)

	return json.NewDecoder(res.Body).Decode(out)
}

func ParseJsonBody(res *http.Response, out interface{}) (err error) {
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if err == nil {
			err = closeErr
		}
	}(res.Body)
	return json.NewDecoder(res.Body).Decode(out)
}
