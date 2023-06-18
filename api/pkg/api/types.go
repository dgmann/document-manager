package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Category struct {
	Id   string `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name,omitempty" json:"name"`
}

func NewCategory(id, name string) *Category {
	return &Category{Id: id, Name: name}
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
	UpdatedAt   time.Time          `bson:"updatedAt"`
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
		"updatedAt":   r.UpdatedAt,
	}
	return json.Marshal(m)
}

func (r *Record) Clone() Record {
	clone := &Record{}
	*clone = *r
	pages := make([]Page, len(r.Pages))
	for i, p := range r.Pages {
		pages[i] = *p.Clone()
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
		UpdatedAt:  time.Now(),
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

type Page struct {
	Id        string    `bson:"id" json:"id"`
	Url       string    `json:"url"`
	Content   string    `bson:"content" json:"content"`
	Format    string    `bson:"format" json:"format"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func NewPage(format string) *Page {
	id := primitive.NewObjectID().Hex()
	return &Page{Id: id, Format: format, UpdatedAt: time.Now()}
}

func (p *Page) Clone() *Page {
	return &Page{
		Id:        p.Id,
		Url:       p.Url,
		Content:   p.Content,
		Format:    p.Format,
		UpdatedAt: p.UpdatedAt,
	}
}

type PageUpdate struct {
	Id      string  `json:"id"`
	Rotate  float64 `json:"rotate,omitempty"`
	Content string  `json:"content,omitempty"`
}

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

func (s Status) IsValid() bool {
	return s == StatusNone || s == StatusInbox || s == StatusEscalated || s == StatusReview || s == StatusOther || s == StatusDone
}

type Type string

const (
	TypeCreated Type = "CREATE"
	TypeUpdated      = "UPDATE"
	TypeDeleted      = "DELETE"
)

type Event struct {
	Type      Type        `json:"type"`
	Topic     Topic       `json:"topic"`
	Timestamp time.Time   `json:"timestamp"`
	Id        string      `json:"id"`
	Data      interface{} `json:"data"`
}

type Topic string

const RecordTopic = "records"

func New(topic Topic, eventType Type, id string, data interface{}) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Id:        id,
		Topic:     topic,
		Data:      data,
	}
}
