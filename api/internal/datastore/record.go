package datastore

import (
	"context"
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
	Query(ctx context.Context, query RecordQuery, queryOption ...*QueryOptions) ([]api.Record, error)
	Create(ctx context.Context, data api.CreateRecord, images []storage.Image, pdfData io.Reader) (*api.Record, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, record api.Record, updateOptions ...UpdateOption) (*api.Record, error)
	UpdatePages(ctx context.Context, id string, updates []api.PageUpdate) (*api.Record, error)
}

type Record struct {
	Id          primitive.ObjectID `bson:"_id"`
	*api.Record `bson:"inline"`
}

type RecordQuery bson.D

type RecordQueryFunc = func(query *RecordQuery)

func NewRecordQuery(queryFuncs ...RecordQueryFunc) RecordQuery {
	query := RecordQuery{}
	for _, queryFunc := range queryFuncs {
		queryFunc(&query)
	}
	return query
}

func WithStatus(status api.Status) RecordQueryFunc {
	return func(query *RecordQuery) {
		if status.IsValid() {
			*query = append(*query, bson.E{Key: "status", Value: status})
		}
	}
}

func WithPatientId(patientId string) RecordQueryFunc {
	return func(query *RecordQuery) {
		*query = append(*query, bson.E{Key: "patientId", Value: patientId})
	}
}

func WithIds(queryIds []string) RecordQueryFunc {
	return func(query *RecordQuery) {
		ids := make([]primitive.ObjectID, len(queryIds))
		for i, id := range queryIds {
			oid, err := primitive.ObjectIDFromHex(id)
			// Ignore invalid ones
			if err != nil {
				continue
			}
			ids[i] = oid
		}
		*query = append(*query, bson.E{Key: "_id", Value: bson.M{"$in": ids}})
	}
}

func WithNoContent() RecordQueryFunc {
	return func(query *RecordQuery) {
		*query = append(*query, bson.E{Key: "pages.content", Value: bson.D{
			{"$not", bson.D{
				{"$nin", []interface{}{bson.TypeNull, ""}},
			}},
		}})
	}
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

type UpdateOption = func(query bson.M)

func IfNotModifiedSince(modified time.Time) UpdateOption {
	return func(query bson.M) {
		query["updatedAt"] = modified
	}
}

func NewPage(format string) *api.Page {
	id := primitive.NewObjectID().Hex()
	return &api.Page{Id: id, Format: format, UpdatedAt: time.Now()}
}
