package mongo

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/api/datastore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	Categories = "categories"
	Records    = "records"
)

type Client struct {
	uri    string
	dbName string
	*mongo.Client
}

func NewClient(config datastore.DatabaseConfig) *Client {
	return &Client{
		uri:    fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port),
		dbName: config.Name,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.uri))
	if err != nil {
		return err
	}
	c.Client = client

	return c.Client.Ping(ctx, readpref.Primary())
}

func (c *Client) Database() *mongo.Database {
	return c.Client.Database(c.dbName)
}

func (c *Client) Records() *mongo.Collection {
	return c.Database().Collection("records")
}

func (c *Client) Categories() *mongo.Collection {
	return c.Database().Collection("categories")
}

func (c *Client) CreateIndexes(ctx context.Context) error {
	_, err := c.Records().Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{"patientId", int32(1)}}})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Check(ctx context.Context) (string, error) {
	err := c.Ping(ctx, readpref.Primary())
	if err != nil {
		return "", err
	}
	return "connected", nil
}
