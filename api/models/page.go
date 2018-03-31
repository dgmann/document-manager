package models

type Page struct {
	Id      string `bson:"id" json:"-"`
	Url     string `bson:"-" json:"url"`
	Content string `json:"content"`
}
