package client

import (
	"github.com/dgmann/document-manager/api/pkg/api"
)

type CategoryClient struct {
	*httpClient
}

func (c *CategoryClient) Create(category *Category) (*api.Category, error) {
	res, err := c.PostJson("categories", category)
	if err != nil {
		return nil, err
	}
	return ToPointer(HandleResponse[api.Category](res))
}

func (c *CategoryClient) Get(id string) (*api.Category, error) {
	res, err := c.GetJson("categories/" + id)
	if err != nil {
		return nil, err
	}
	return ToPointer(HandleResponse[api.Category](res))
}

func (c *CategoryClient) All() ([]api.Category, error) {
	res, err := c.GetJson("categories")
	if err != nil {
		return nil, err
	}
	return HandleResponse[[]api.Category](res)
}
