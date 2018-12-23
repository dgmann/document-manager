package category

import (
	"context"
	"github.com/dgmann/document-manager/api/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	All(ctx context.Context) ([]models.Category, error)
	Find(ctx context.Context, id string) (*models.Category, error)
	FindByPatient(ctx context.Context, id string) ([]models.Category, error)
	Add(ctx context.Context, id, category string) error
}

type DatabaseRepository struct {
	categories *mongo.Collection
	records    *mongo.Collection
}

func NewDatabaseRepository(categories, records *mongo.Collection) *DatabaseRepository {
	return &DatabaseRepository{categories: categories, records: records}
}

func (c *DatabaseRepository) All(ctx context.Context) ([]models.Category, error) {
	cursor, err := c.categories.Find(ctx, bson.M{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return castToSlice(ctx, cursor)
}

func (c *DatabaseRepository) Find(ctx context.Context, id string) (*models.Category, error) {
	var category models.Category

	if err := c.categories.FindOne(ctx, bson.M{"_id": id}).Decode(&category); err != nil {
		log.WithField("error", err).Panic("Cannot find record")
		return nil, err
	}
	return &category, nil
}

func (c *DatabaseRepository) FindByPatient(ctx context.Context, id string) ([]models.Category, error) {
	ids, err := c.records.Distinct(ctx, "category", bson.M{"patientId": id})
	if err != nil {
		log.WithField("error", err).Panic("Cannot find categories by patient id")
		return nil, err
	}

	cursor, err := c.categories.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		log.WithField("error", err).Panic("Cannot resolve category ids")
		return nil, err
	}
	return castToSlice(ctx, cursor)
}

func (c *DatabaseRepository) Add(ctx context.Context, id, category string) error {
	_, err := c.categories.InsertOne(ctx, models.NewCategory(id, category))
	if err != nil {
		return err
	}
	return nil
}
