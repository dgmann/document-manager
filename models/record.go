package models

import "time"

type Record struct {
	Id        int64     `jsonapi:"primary,records"`
	Date      time.Time `jsonapi:"attr,date,iso8601"`
	Comment   string    `jsonapi:"attr,comment"`
	Sender    string    `jsonapi:"attr,sender" form:"user" binding:"required"`
	Pages     []Page    `jsonapi:"attr,pages"`
	Processed bool      `jsonapi:"attr,processsed"`
	Escalated bool      `jsonapi:"attr,escalated"`
}
