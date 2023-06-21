package datastore

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/api/internal/storage"
	"github.com/dgmann/document-manager/api/pkg/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"time"
)

type RecordService interface {
	All(ctx context.Context) ([]api.Record, error)
	Find(ctx context.Context, id string) (*api.Record, error)
	Query(ctx context.Context, query *RecordQuery, queryOption ...*QueryOptions) ([]api.Record, error)
	Create(ctx context.Context, data api.CreateRecord, images []storage.Image, pdfData io.Reader) (*api.Record, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, record api.Record) (*api.Record, error)
	UpdatePages(ctx context.Context, id string, updates []api.PageUpdate) (*api.Record, error)
}

type Record struct {
	Id          primitive.ObjectID `bson:"_id"`
	*api.Record `bson:"inline"`
}

type RecordQuery struct {
	Status    api.Status
	PatientId *string
	Ids       []string
}

func NewRecordQuery() *RecordQuery {
	return &RecordQuery{}
}

func (query *RecordQuery) SetStatus(status api.Status) *RecordQuery {
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

func NewPage(format string) *api.Page {
	id := primitive.NewObjectID().Hex()
	return &api.Page{Id: id, Format: format, UpdatedAt: time.Now()}
}
