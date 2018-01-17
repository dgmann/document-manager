package models

import "time"

type Record struct {
	Id        string    `jsonapi:"primary,records"`
	Date      time.Time `jsonapi:"attr,date,iso8601"`
	Comment   string    `jsonapi:"attr,comment"`
	Sender    string    `jsonapi:"attr,sender"`
	Pages     []Page    `jsonapi:"attr,pages"`
	Processed bool      `jsonapi:"attr,processsed"`
	Escalated bool      `jsonapi:"attr,escalated"`
}

type Page struct {
	Id      int    `json:"-"`
	Index   int    `json:"index"`
	Url     string `json:"url"`
	Content string `json:"content"`
}
