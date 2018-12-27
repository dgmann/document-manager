package app

import (
	"context"
)

type Category struct {
	Id   string `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name,omitempty" json:"name"`
}

func NewCategory(id, name string) *Category {
	return &Category{Id: id, Name: name}
}

type CategoryService interface {
	All(ctx context.Context) ([]Category, error)
	Find(ctx context.Context, id string) (*Category, error)
	FindByPatient(ctx context.Context, id string) ([]Category, error)
	Add(ctx context.Context, id, category string) error
}
