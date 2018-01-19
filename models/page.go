package models

type Page struct {
	Id      int    `json:"-"`
	Index   int    `json:"index"`
	Url     string `json:"url"`
	Content string `json:"content"`
}
