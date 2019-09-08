package app

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

type Status string

const (
	StatusNone      Status = ""
	StatusInbox     Status = "inbox"
	StatusEscalated Status = "escalated"
	StatusReview    Status = "review"
	StatusOther     Status = "other"
	StatusDone      Status = "done"
)

func (s Status) IsNone() bool {
	return s == StatusNone
}

type RecordService interface {
	All(ctx context.Context) ([]Record, error)
	Find(ctx context.Context, id string) (*Record, error)
	Query(ctx context.Context, query map[string]interface{}) ([]Record, error)
	Create(ctx context.Context, data CreateRecord, images []Image, pdfData io.Reader) (*Record, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, record Record) (*Record, error)
	UpdatePages(ctx context.Context, id string, updates []PageUpdate) (*Record, error)
}

type Record struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date        *time.Time         `bson:"date,omitempty" json:"date"`
	ReceivedAt  time.Time          `bson:"receivedAt,omitempty" json:"receivetAt"`
	PatientId   *string            `bson:"patientId,omitempty" json:"patientId"`
	Comment     *string            `bson:"comment,omitempty" json:"comment"`
	Sender      string             `bson:"sender,omitempty" json:"sender" form:"user" binding:"required"`
	Category    *string            `bson:"category,omitempty" json:"category"`
	Tags        *[]string          `bson:"tags,omitempty" json:"tags"`
	Pages       []Page             `bson:"pages,omitempty" json:"pages"`
	Status      *Status            `bson:"status,omitempty" json:"status"`
	ArchivedPDF string
}

func (r *Record) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"id":          r.Id,
		"date":        formatTime(r.Date),
		"receivedAt":  formatTime(&r.ReceivedAt),
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

func (r *Record) Clone() Record {
	clone := &Record{}
	*clone = *r
	pages := make([]Page, len(r.Pages))
	for i, p := range r.Pages {
		page := &Page{}
		if err := copier.Copy(page, p); err != nil {
			logrus.WithError(err).Warn("error cloning page")
		}
		pages[i] = *page
	}
	clone.Pages = pages

	return *clone
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

func formatTime(toFormat *time.Time) *string {
	if toFormat == nil {
		return nil
	}
	res := toFormat.Format(time.RFC3339)
	return &res
}

type CreateRecord struct {
	Id         *primitive.ObjectID `form:"id"`
	Date       time.Time           `form:"date" time_format:"2006-01-02T15:04:05Z07:00"`
	ReceivedAt time.Time           `form:"receivedAt" time_format:"2006-01-02T15:04:05Z07:00"`
	Sender     string              `form:"sender"`
	Comment    *string             `form:"comment"`
	PatientId  *string             `form:"patientId"`
	Tags       []string            `form:"tags"`
	Status     Status              `form:"status"`
	Category   *string             `form:"category"`
	Pages      []Page
}

func NewRecord(data CreateRecord) *Record {
	id := primitive.NewObjectID()
	if data.Id != nil {
		id = *data.Id
	}
	record := &Record{
		Id:         id,
		Date:       nil,
		ReceivedAt: time.Now(),
		Comment:    data.Comment,
		PatientId:  data.PatientId,
		Sender:     data.Sender,
		Tags:       &data.Tags,
		Pages:      data.Pages,
		Status:     &data.Status,
		Category:   data.Category,
	}
	if !data.Date.IsZero() {
		record.Date = &data.Date
	}
	if !data.ReceivedAt.IsZero() {
		record.ReceivedAt = data.ReceivedAt
	}
	if len(*record.Tags) == 0 {
		record.Tags = &[]string{}
	}
	if len(record.Pages) == 0 {
		record.Pages = []Page{}
	}
	if record.Status.IsNone() {
		status := StatusInbox
		record.Status = &status
	}
	return record
}
