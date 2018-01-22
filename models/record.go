package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Record struct {
	Id        string        `jsonapi:"primary,records"`
	Primary   bson.ObjectId `bson:"_id,omitempty"` //Primary key for mongodb. Not serialized
	Date      time.Time     `bson:"date" jsonapi:"attr,date,iso8601"`
	Comment   string        `bson:"comment" jsonapi:"attr,comment"`
	Sender    string        `bson:"sender" jsonapi:"attr,sender" form:"user" binding:"required"`
	Tags      []string      `bson:"tags" jsonapi:"attr,tags"`
	Pages     []Page        `bson:"pages" jsonapi:"attr,pages"`
	Processed *bool         `bson:"processed" jsonapi:"attr,processed"`
	Escalated *bool         `bson:"escalated" jsonapi:"attr,escalated"`
}

func NewRecord(id bson.ObjectId, sender string) *Record {
	return &Record{Primary: id,
		Date:      time.Now(),
		Comment:   "",
		Sender:    sender,
		Tags:      []string{},
		Pages:     []Page{},
		Processed: newFalse(),
		Escalated: newFalse(),
	}
}

func newFalse() *bool {
	b := false
	return &b
}
