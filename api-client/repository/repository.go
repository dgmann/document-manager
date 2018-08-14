package repository

import (
	"net/http"
	"bytes"
	"encoding/json"
)

type Repository struct {
	url string
}

func NewRepository(url string) *Repository {
	return &Repository{url}
}

func (r *Repository) Create(path string, model interface{}) error {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(model)

	res, err := http.Post(r.url+path, "application/json", body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
