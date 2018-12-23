package tag

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type Repository interface {
	All(ctx context.Context) ([]string, error)
	ByPatient(ctx context.Context, id string) ([]string, error)
}

type DatabaseRepository struct {
	records *mongo.Collection
}

func NewDatabaseRepository(records *mongo.Collection) *DatabaseRepository {
	return &DatabaseRepository{records: records}
}

func (t *DatabaseRepository) All(ctx context.Context) ([]string, error) {
	return t.query(ctx, bson.M{})
}

func (t *DatabaseRepository) ByPatient(ctx context.Context, id string) ([]string, error) {
	return t.query(ctx, bson.M{"patientId": id})
}

func (t *DatabaseRepository) query(ctx context.Context, query interface{}) ([]string, error) {
	var tags []string
	res, err := t.records.Distinct(ctx, "tags", bson.M{})
	if err != nil {
		return tags, err
	}
	for _, tag := range res {
		tags = append(tags, tag.(string))
	}
	return tags, err
}
