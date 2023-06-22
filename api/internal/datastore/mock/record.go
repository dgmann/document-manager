package mock

import (
	"context"
	"github.com/dgmann/document-manager/api/internal/datastore"
	"github.com/dgmann/document-manager/api/internal/storage"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/stretchr/testify/mock"
	"io"
)

type RecordService struct {
	mock.Mock
}

func NewRecordService() *RecordService {
	return &RecordService{}
}

func (mock *RecordService) All(ctx context.Context) ([]api.Record, error) {
	args := mock.Called(ctx)
	return args.Get(0).([]api.Record), args.Error(1)
}
func (mock *RecordService) Find(ctx context.Context, id string) (*api.Record, error) {
	args := mock.Called(ctx, id)
	return args.Get(0).(*api.Record), args.Error(1)
}
func (mock *RecordService) Query(ctx context.Context, query datastore.RecordQuery, options ...*datastore.QueryOptions) ([]api.Record, error) {
	o := make([]interface{}, len(options)+2)
	o[0] = ctx
	o[1] = query
	for i, opt := range options {
		o[i+2] = opt
	}
	args := mock.Called(o...)
	return args.Get(0).([]api.Record), args.Error(1)
}
func (mock *RecordService) Create(ctx context.Context, data api.CreateRecord, images []storage.Image, pdfData io.Reader) (*api.Record, error) {
	args := mock.Called(ctx, data, images, pdfData)
	return args.Get(0).(*api.Record), args.Error(1)
}
func (mock *RecordService) Delete(ctx context.Context, id string) error {
	args := mock.Called(ctx, id)
	return args.Error(0)
}
func (mock *RecordService) Update(ctx context.Context, id string, record api.Record) (*api.Record, error) {
	args := mock.Called(ctx, id, record)
	return args.Get(0).(*api.Record), args.Error(1)
}
func (mock *RecordService) UpdatePages(ctx context.Context, id string, updates []api.PageUpdate) (*api.Record, error) {
	args := mock.Called(ctx, id, updates)
	return args.Get(0).(*api.Record), args.Error(1)
}
