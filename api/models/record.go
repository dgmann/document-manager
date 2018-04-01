package models

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"net/url"
	"time"
)

const (
	ActionEscalated = "escalated"
	ActionReview    = "review"
	ActionOther     = "other"
)

type Record struct {
	Id             string        `json:"id"`
	Primary        bson.ObjectId `bson:"_id,omitempty" json:"-"` //Primary key for mongodb. Not serialized
	Date           *time.Time    `bson:"date,omitempty" json:"date"`
	ReceivedAt     time.Time     `bson:"receivedAt,omitempty" json:"receivetAt"`
	PatientId      *string       `bson:"patientId,omitempty" json:"patientId"`
	Comment        *string       `bson:"comment,omitempty" json:"comment"`
	Sender         string        `bson:"sender,omitempty" json:"sender" form:"user" binding:"required"`
	Category       *string       `bson:"category,omitempty" json:"category"`
	Tags           []string      `bson:"tags,omitempty" json:"tags"`
	Pages          []Page        `bson:"pages,omitempty" json:"pages"`
	RequiredAction *string       `bson:"requiredAction,omitempty" json:"requiredAction"`
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
		"id":             r.Id,
		"date":           r.Date,
		"receivedAt":     r.ReceivedAt,
		"patientId":      toString(r.PatientId),
		"comment":        toString(r.Comment),
		"sender":         r.Sender,
		"category":       r.Category,
		"tags":           r.Tags,
		"pages":          r.Pages,
		"requiredAction": toString(r.RequiredAction),
	}
	return json.Marshal(m)
}

func toString(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

func NewRecord(id bson.ObjectId, sender string) *Record {
	return &Record{Primary: id,
		Date:           nil,
		ReceivedAt:     time.Now(),
		Comment:        nil,
		PatientId:      nil,
		Sender:         sender,
		Tags:           []string{},
		Pages:          []Page{},
		RequiredAction: nil,
	}
}

func (r *Record) SetPages(ids []string) {
	for _, id := range ids {
		r.Pages = append(r.Pages, Page{Id: id})
	}
}