package api

import (
	"encoding/json"
	"time"
)

type MatchType string

const (
	MatchTypeNone  = ""
	MatchTypeRegex = "regex"
	MatchTypeExact = "exact"
	MatchTypeAll   = "all"
	MatchTypeAny   = "any"
)

type MatchConfig struct {
	Type       MatchType `bson:"type,omitempty" json:"type"`
	Expression string    `bson:"expression,omitempty" json:"expression"`
}

type Category struct {
	Id    string      `bson:"_id,omitempty" json:"id"`
	Name  string      `bson:"name,omitempty" json:"name"`
	Match MatchConfig `bson:"match,omitempty" json:"match"`
}

func NewCategory(id, name string, matchConfig MatchConfig) *Category {
	return &Category{Id: id, Name: name, Match: matchConfig}
}

type Record struct {
	Id          string     `bson:"_id,omitempty" json:"id"`
	Date        *time.Time `bson:"date,omitempty" json:"date"`
	ReceivedAt  time.Time  `bson:"receivedAt,omitempty" json:"receivedAt"`
	PatientId   *string    `bson:"patientId,omitempty" json:"patientId"`
	Comment     *string    `bson:"comment,omitempty" json:"comment"`
	Sender      string     `bson:"sender,omitempty" json:"sender" form:"user" binding:"required"`
	Category    *string    `bson:"category,omitempty" json:"category"`
	Tags        *[]string  `bson:"tags,omitempty" json:"tags"`
	Pages       []Page     `bson:"pages,omitempty" json:"pages"`
	Status      *Status    `bson:"status,omitempty" json:"status"`
	UpdatedAt   time.Time  `bson:"updatedAt"`
	ArchivedPDF string     `json:"archivedPDF"`
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
		"status":      r.Status.String(),
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

func formatTime(toFormat *time.Time) *string {
	if toFormat == nil {
		return nil
	}
	res := toFormat.Format(time.RFC3339)
	return &res
}

type CreateRecord struct {
	Date       time.Time `form:"date" time_format:"2006-01-02T15:04:05Z07:00"`
	ReceivedAt time.Time `form:"receivedAt" time_format:"2006-01-02T15:04:05Z07:00"`
	Sender     string    `form:"sender"`
	Comment    *string   `form:"comment"`
	PatientId  *string   `form:"patientId"`
	Tags       []string  `form:"tags"`
	Status     Status    `form:"status"`
	Category   *string   `form:"category"`
	Pages      []Page
}

type Page struct {
	Id        string    `bson:"id" json:"id"`
	Url       string    `json:"url"`
	Content   *string   `bson:"content,omitempty" json:"content"`
	Format    string    `bson:"format" json:"format"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
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
	Content *string `json:"content,omitempty"`
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

func (s *Status) IsNone() bool {
	status := s.Status()
	return status == StatusNone
}

func (s *Status) IsValid() bool {
	status := s.Status()
	return status == StatusNone || status == StatusInbox || status == StatusEscalated || status == StatusReview || status == StatusOther || status == StatusDone
}

func (s *Status) String() string {
	return string(s.Status())
}

func (s *Status) Status() Status {
	if s == nil {
		return StatusNone
	}
	return *s
}

type EventType string

type Topic string

const RecordTopic = "records"

const (
	EventTypeCreated EventType = "CREATE"
	EventTypeUpdated           = "UPDATE"
	EventTypeDeleted           = "DELETE"
)

type Event[T any] struct {
	Type      EventType `json:"type"`
	Topic     Topic     `json:"topic"`
	Timestamp time.Time `json:"timestamp"`
	Id        string    `json:"id"`
	Data      T         `json:"data"`
}

func NewEvent[T any](topic Topic, eventType EventType, id string, data T) Event[T] {
	return Event[T]{
		Type:      eventType,
		Timestamp: time.Now(),
		Id:        id,
		Topic:     topic,
		Data:      data,
	}
}
