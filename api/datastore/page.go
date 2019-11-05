package datastore

import "github.com/mongodb/mongo-go-driver/bson/primitive"

type Page struct {
	Id      string `bson:"id" json:"id"`
	Url     string `json:"url"`
	Content string `bson:"content" json:"content"`
	Format  string `bson:"format" json:"format"`
}

func NewPage(format string) *Page {
	id := primitive.NewObjectID().Hex()
	return &Page{Id: id, Format: format}
}

type PageUpdate struct {
	Id     string  `json:"id"`
	Rotate float64 `json:"rotate"`
}
