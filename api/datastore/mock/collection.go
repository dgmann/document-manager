package mock

import (
	"context"
	"github.com/dgmann/document-manager/api/datastore"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	mock.Mock
}

func NewCollection() *Collection {
	return &Collection{}
}

func (m *Collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (datastore.Cursor, error) {
	params := []interface{}{ctx, filter}
	for _, opt := range opts {
		params = append(params, opt)
	}
	args := m.Called(params...)
	return args.Get(0).(datastore.Cursor), args.Error(1)
}

func (m *Collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) datastore.Decodable {
	params := []interface{}{ctx, filter}
	for _, opt := range opts {
		params = append(params, opt)
	}
	args := m.Called(params...)
	return args.Get(0).(datastore.Decodable)
}

func (m *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	params := []interface{}{ctx, document}
	for _, opt := range opts {
		params = append(params, opt)
	}
	args := m.Called(params...)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *Collection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	params := []interface{}{ctx, fieldName, filter}
	for _, opt := range opts {
		params = append(params, opt)
	}
	args := m.Called(params...)
	return args.Get(0).([]interface{}), args.Error(1)
}
