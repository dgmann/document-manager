package models

import (
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"github.com/jinzhu/copier"
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
	Category    *string       `bson:"category,omitempty" json:"category"`
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
		"category":    r.Category,
		"tags":        r.Tags,
		"pages":       r.Pages,
		"status":      statusToString(r.Status),
		"archivedPDF": r.ArchivedPDF,
	}
	return json.Marshal(m)
}

func (r *Record) Clone() *Record {
	clone := &Record{}
	*clone = *r
	pages := make([]*Page, len(r.Pages))
	for i, p := range r.Pages {
		page := &Page{}
		copier.Copy(page, p)
		pages[i] = page
	}
	clone.Pages = pages

	return clone
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

type CreateRecord struct {
	Date       time.Time `form:"date" time_format:"2006-01-02T15:04:05Z07:00"`
	ReceivedAt time.Time `form:"receivedAt" time_format:"2006-01-02T15:04:05Z07:00"`
	Sender     string    `form:"sender"`
	Comment    *string   `form:"comment"`
	PatientId  *string   `form:"patientId"`
	Tags       []string  `form:"tags"`
	Status     *Status   `form:"status"`
	Category   *string   `form:"category"`
}

func NewRecord(data CreateRecord) *Record {
	record := &Record{
		Id:         bson.NewObjectId(),
		Date:       nil,
		ReceivedAt: time.Now(),
		Comment:    data.Comment,
		PatientId:  data.PatientId,
		Sender:     data.Sender,
		Tags:       data.Tags,
		Pages:      []*Page{},
		Status:     data.Status,
		Category:   data.Category,
	}
	if !data.Date.IsZero() {
		record.Date = &data.Date
	}
	if !data.ReceivedAt.IsZero() {
		record.ReceivedAt = data.ReceivedAt
	}
	if record.Tags == nil {
		record.Tags = []string{}
	}
	if record.Status == nil {
		status := StatusInbox
		record.Status = &status
	}
	return record
}
