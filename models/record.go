package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Record struct {
	Id         string        `jsonapi:"primary,records"`
	Primary    bson.ObjectId `bson:"_id,omitempty"` //Primary key for mongodb. Not serialized
	Date       time.Time     `bson:"date,omitempty" jsonapi:"attr,date,iso8601"`
	ReceivedAt time.Time     `bson:"receivedAt,omitempty" jsonapi:"attr,receivetAt,iso8601"`
	PatientId  string        `bson:"patientId,omitempty" jsonapi:"attr,patientId"`
	Comment    string        `bson:"comment,omitempty" jsonapi:"attr,comment"`
	Sender     string        `bson:"sender,omitempty" jsonapi:"attr,sender" form:"user" binding:"required"`
	Tags       []string      `bson:"tags,omitempty" jsonapi:"attr,tags"`
	Pages      []Page        `bson:"pages,omitempty" jsonapi:"attr,pages"`
	Processed  *bool         `bson:"processed,omitempty" jsonapi:"attr,processed"`
	Escalated  *bool         `bson:"escalated,omitempty" jsonapi:"attr,escalated"`
}

func NewRecord(id bson.ObjectId, sender string) *Record {
	return &Record{Primary: id,
		Date:       time.Now(),
		ReceivedAt: time.Now(),
		Comment:    "",
		PatientId:  "",
		Sender:     sender,
		Tags:       []string{},
		Pages:      []Page{},
		Processed:  newFalse(),
		Escalated:  newFalse(),
	}
}

func newFalse() *bool {
	b := false
	return &b
}
