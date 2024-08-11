//go:build test

// This file in only required to not require any tesseract libraries during testing

package tesseract

import (
	"ocr/internal/ocr"

	"github.com/dgmann/document-manager/pkg/api"
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
