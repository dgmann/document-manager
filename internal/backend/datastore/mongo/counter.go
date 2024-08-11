package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Counter struct {
	client *Client
}

func NewCounter(client *Client) *Counter {
	return &Counter{client: client}
}

func (c *Counter) Count(ctx context.Context, resource string) (int64, error) {
	if resource == CollectionCategories || resource == CollectionRecords {
		return c.client.Database().Collection(resource).CountDocuments(ctx, bson.D{})
	} else if resource == "patients" {
		elements, err := c.client.Database().Collection(CollectionRecords).Aggregate(ctx, bson.A{bson.D{{"$group", bson.M{"_id": "$patientId"}}}, bson.D{{"$group", bson.D{{"_id", 1}, {"count", bson.M{"$sum": 1}}}}}})
		if err != nil {
			return 0, err
		}
		var s []primitive.D
		err = elements.All(ctx, &s)
		if len(s) != 1 {
			return 0, errors.New("unexpected result")
		}
		return int64(s[0].Map()["count"].(int32)), err
	}
	return 0, nil
}
