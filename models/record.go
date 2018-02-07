package models

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"net/url"
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

func (r *Record) SetURL(url *url.URL) {
	if r.Tags == nil {
		r.Tags = []string{}
	}
	if r.Pages == nil {
		r.Pages = []Page{}
	}

	for i := range r.Pages {
		r.Pages[i].Url = fmt.Sprintf("%s/records/%s/images/%s", url.String(), r.Id, r.Pages[i].Id)
	}
}

func (r *Record) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"id":         r.Id,
		"date":       r.Date,
		"receivedAt": r.ReceivedAt,
		"patientId":  r.PatientId,
		"comment":    r.Comment,
		"sender":     r.Sender,
		"tags":       r.Tags,
		"pages":      r.Pages,
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
