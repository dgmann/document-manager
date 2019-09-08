package mock

import (
	"context"
	"github.com/dgmann/document-manager/api/app"
	"github.com/stretchr/testify/mock"
	"io"
)

type RecordService struct {
	mock.Mock
}

func (mock *RecordService) All(ctx context.Context) ([]app.Record, error) {
	args := mock.Called(ctx)
	return args.Get(0).([]app.Record), args.Error(1)
}
func (mock *RecordService) Find(ctx context.Context, id string) (*app.Record, error) {
	args := mock.Called(ctx, id)
	return args.Get(0).(*app.Record), args.Error(1)
}
func (mock *RecordService) Query(ctx context.Context, query map[string]interface{}) ([]app.Record, error) {
	args := mock.Called(ctx, query)
	return args.Get(0).([]app.Record), args.Error(1)
}
func (mock *RecordService) Create(ctx context.Context, data app.CreateRecord, images []app.Image, pdfData io.Reader) (*app.Record, error) {
	args := mock.Called(ctx, data, images, pdfData)
	return args.Get(0).(*app.Record), args.Error(1)
}
func (mock *RecordService) Delete(ctx context.Context, id string) error {
	args := mock.Called(ctx, id)
	return args.Error(0)
}
func (mock *RecordService) Update(ctx context.Context, id string, record app.Record) (*app.Record, error) {
	args := mock.Called(ctx, id, record)
	return args.Get(0).(*app.Record), args.Error(1)
}
func (mock *RecordService) UpdatePages(ctx context.Context, id string, updates []app.PageUpdate) (*app.Record, error) {
	args := mock.Called(ctx, id, updates)
	return args.Get(0).(*app.Record), args.Error(1)
}
