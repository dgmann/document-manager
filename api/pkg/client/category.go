package client

import (
	"fmt"
	"net/http"
)

type CategoryClient struct {
	*httpClient
}

func (c *CategoryClient) Create(category *Category) error {
	res, err := c.PostJson("categories", category)
	if err != nil {
		return err
	}
	var body map[string]interface{}
	if err := ParseJsonBody(res, &body); err != nil {
		return err
	}
	if res.StatusCode == http.StatusCreated {
		return nil
	}
	if res.StatusCode == http.StatusConflict {
		return fmt.Errorf("category %s could not be created: %w", category.Id, ErrAlreadyExists)
	}
	return fmt.Errorf("status code: %d, body: %+v", res.StatusCode, body)
}
