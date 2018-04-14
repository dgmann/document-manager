package models

type Page struct {
	Id      string `bson:"id" json:"-"`
	Url     string `json:"url"`
	Content string `bson:"content" json:"content"`
	Format  string `bson:"format" json:"format"`
}

func NewPage(id, format string) *Page {
	return &Page{Id: id, Format: format}
}
