package models

import (
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"time"
)

type Record struct {
	Id         string        `json:"id"`
	Primary    bson.ObjectId `bson:"_id,omitempty" json:"-"` //Primary key for mongodb. Not serialized
	Date       time.Time     `bson:"date,omitempty" json:"date"`
	ReceivedAt time.Time     `bson:"receivedAt,omitempty" json:"receivetAt"`
	PatientId  string        `bson:"patientId,omitempty" json:"patientId"`
	Comment    string        `bson:"comment,omitempty" json:"comment"`
	Sender     string        `bson:"sender,omitempty" json:"sender" form:"user" binding:"required"`
	Tags       []string      `bson:"tags,omitempty" json:"tags"`
	Pages      []Page        `bson:"pages,omitempty" json:"pages"`
	Processed  *bool         `bson:"processed,omitempty" json:"processed"`
	Escalated  *bool         `bson:"escalated,omitempty" json:"escalated"`
}

func (p *Record) MarshalJSON() ([]byte, error) {
	tags := p.Tags
	if tags == nil {
		tags = []string{}
	}
	pages := p.Pages
	if pages == nil {
		pages = []Page{}
	}
	m := map[string]interface{}{
		"id":         p.Id,
		"date":       p.Date,
		"receivedAt": p.ReceivedAt,
		"patientId":  p.PatientId,
		"comment":    p.Comment,
		"sender":     p.Sender,
		"tags":       tags,
		"pages":      pages,
		"processed":  p.Processed,
		"escalated":  p.Escalated,
	}
	return json.Marshal(m)
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
