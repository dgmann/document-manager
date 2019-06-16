package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	url string
}

func NewClient(url string) *Client {
	return &Client{url}
}

func (r *Client) Create(path string, model interface{}) error {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(model); err != nil {
		return err
	}

	res, err := http.Post(r.url+path, "application/json", body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
