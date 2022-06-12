package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/dgmann/document-manager/api/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s Status) IsValid() bool {
	return s == StatusNone || s == StatusInbox || s == StatusEscalated || s == StatusReview || s == StatusOther || s == StatusDone
}

type RecordService interface {
	All(ctx context.Context) ([]Record, error)
	Find(ctx context.Context, id string) (*Record, error)
	Query(ctx context.Context, query *RecordQuery, queryOption ...*QueryOptions) ([]Record, error)
	Create(ctx context.Context, data CreateRecord, images []storage.Image, pdfData io.Reader) (*Record, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, record Record) (*Record, error)
	UpdatePages(ctx context.Context, id string, updates []PageUpdate) (*Record, error)
}

type RecordQuery struct {
	Status    Status
	PatientId *string
	Ids       []string
}

func NewRecordQuery() *RecordQuery {
	return &RecordQuery{}
}

func (query *RecordQuery) SetStatus(status Status) *RecordQuery {
	if status.IsValid() {
		query.Status = status
	}
	return query
}

func (query *RecordQuery) SetPatientId(patientId string) *RecordQuery {
	query.PatientId = &patientId
	return query
}

func (query *RecordQuery) SetIds(ids []string) *RecordQuery {
	query.Ids = ids
	return query
}

func (query *RecordQuery) ToMap() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if !query.Status.IsNone() {
		result["status"] = query.Status
	}
	if query.Ids != nil {
		ids := make([]primitive.ObjectID, len(query.Ids))
		for i, id := range query.Ids {
			oid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return nil, fmt.Errorf("invalid id %s: %w", id, err)
			}
			ids[i] = oid
		}
		result["_id"] = bson.M{"$in": ids}
	}
	if query.PatientId != nil {
		result["patientId"] = *query.PatientId
	}
	return result, nil
}

type QueryOptions struct {
	Sort  map[string]int
	Skip  int64
	Limit int64
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{Sort: map[string]int{}}
}

func (options *QueryOptions) SetSort(key string) *QueryOptions {
	if key == "" {
		return options
	}
	direction := 1
	if key[0] == '-' {
		direction = -1
		key = key[1:]
	}

	options.Sort[key] = direction
	return options
}

func (options *QueryOptions) SetSkip(value int64) *QueryOptions {
	options.Skip = value
	return options
}

func (options *QueryOptions) SetLimit(value int64) *QueryOptions {
	options.Limit = value
	return options
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
	Id     string  `json:"id"`
	Rotate float64 `json:"rotate"`
}
