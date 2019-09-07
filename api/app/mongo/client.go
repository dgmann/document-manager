package mongo

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"time"
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

func NewClient(dbHost string, dbName string) *Client {
	return &Client{uri: "mongodb://" + dbHost, dbName: dbName}
}

func (c *Client) Connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, c.uri)
	if err != nil {
		return err
	}
	c.Client = client

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
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
