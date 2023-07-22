//go:build test

// This file in only required to not require any tesseract libraries during testing

package tesseract

import (
	"github.com/dgmann/document-manager/api/pkg/api"
	"ocr/internal/ocr"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) CheckOrientation(pages []ocr.PageWithContent) ([]api.PageUpdate, error) {
	return nil, nil
}

func (c *Client) ExtractText(pages []ocr.PageWithContent) ([]api.PageUpdate, error) {
	return nil, nil
}
