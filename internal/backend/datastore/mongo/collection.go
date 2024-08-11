package mongo

import (
	"context"

	"github.com/dgmann/document-manager/internal/backend/datastore"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type finder interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (datastore.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) datastore.Decodable
}

type oneInserter interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

type oneDeleter interface {
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

type distinctFinder interface {
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
}

type oneFinderUpdater interface {
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) datastore.Decodable
}

type Collection struct {
	*mongo.Collection
}

func NewCollection(collection *mongo.Collection) *Collection {
	return &Collection{collection}
}

func (c *Collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (datastore.Cursor, error) {
	return c.Collection.Find(ctx, filter, opts...)
}

func (c *Collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) datastore.Decodable {
	return c.Collection.FindOne(ctx, filter, opts...)
}

func (c *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.Collection.InsertOne(ctx, document, opts...)
}

func (c *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.Collection.DeleteOne(ctx, filter, opts...)
}

func (c *Collection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	return c.Collection.Distinct(ctx, fieldName, filter, opts...)
}

func (c *Collection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) datastore.Decodable {
	return c.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
}
