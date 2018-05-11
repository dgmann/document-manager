package models

import (
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"time"
)

type Status string

const (
	StatusInbox     Status = "inbox"
	StatusEscalated Status = "escalated"
	StatusReview    Status = "review"
	StatusOther     Status = "other"
	StatusDone      Status = "done"
)

type Record struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Date        *time.Time    `bson:"date,omitempty" json:"date"`
	ReceivedAt  time.Time     `bson:"receivedAt,omitempty" json:"receivetAt"`
	PatientId   *string       `bson:"patientId,omitempty" json:"patientId"`
	Comment     *string       `bson:"comment,omitempty" json:"comment"`
	Sender      string        `bson:"sender,omitempty" json:"sender" form:"user" binding:"required"`
	CategoryId  bson.ObjectId `bson:"categoryId,omitempty" json:"categoryId"`
	Tags        []string      `bson:"tags,omitempty" json:"tags"`
	Pages       []*Page       `bson:"pages,omitempty" json:"pages"`
	Status      *Status       `bson:"status,omitempty" json:"status"`
	ArchivedPDF string
}

func (r *Record) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"id":          r.Id,
		"date":        r.Date,
		"receivedAt":  r.ReceivedAt,
		"patientId":   toString(r.PatientId),
		"comment":     toString(r.Comment),
		"sender":      r.Sender,
		"categoryId":  r.CategoryId,
		"tags":        r.Tags,
		"pages":       r.Pages,
		"status":      statusToString(r.Status),
		"archivedPDF": r.ArchivedPDF,
	}
	return json.Marshal(m)
}

func toString(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

func statusToString(val *Status) string {
	if val == nil {
		return ""
	}
	return string(*val)
}

func NewRecord(sender string) *Record {
	status := StatusInbox
	return &Record{
		Id:         bson.NewObjectId(),
		Date:       nil,
		ReceivedAt: time.Now(),
		Comment:    nil,
		PatientId:  nil,
		Sender:     sender,
		Tags:       []string{},
		Pages:      []*Page{},
		Status:     &status,
	}
}
