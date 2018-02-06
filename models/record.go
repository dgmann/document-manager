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
	Pages      []Page        `json:"pages"`
	Processed  *bool         `bson:"processed,omitempty" json:"processed"`
	Escalated  *bool         `bson:"escalated,omitempty" json:"escalated"`
}

func (r *Record) MarshalJSON() ([]byte, error) {
	tags := r.Tags
	if tags == nil {
		tags = []string{}
	}
	pages := r.Pages
	if pages == nil {
		pages = []Page{}
	}
	for i := range pages {

		pages[i].Url = "/records/" + r.Id + "/images/" + pages[i].Id
	}
	m := map[string]interface{}{
		"id":         r.Id,
		"date":       r.Date,
		"receivedAt": r.ReceivedAt,
		"patientId":  r.PatientId,
		"comment":    r.Comment,
		"sender":     r.Sender,
		"tags":       tags,
		"pages":      pages,
		"processed":  r.Processed,
		"escalated":  r.Escalated,
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

func (r *Record) SetPages(ids []string) {
	for _, id := range ids {
		r.Pages = append(r.Pages, Page{Id: id})
	}
}
