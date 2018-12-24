package category

import (
	"context"
	"github.com/dgmann/document-manager/api/models"
	"github.com/dgmann/document-manager/api/repositories/database"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	All(ctx context.Context) ([]models.Category, error)
	Find(ctx context.Context, id string) (*models.Category, error)
	FindByPatient(ctx context.Context, id string) ([]models.Category, error)
	Add(ctx context.Context, id, category string) error
}

type DatabaseRepository struct {
	categories collection
	records    distinctFinder
	decoder    database.Decoder
}

func NewDatabaseRepository(categories collection, records distinctFinder) *DatabaseRepository {
	return &DatabaseRepository{categories: categories, records: records, decoder: database.NewDefaultDecoder()}
}

func (c *DatabaseRepository) All(ctx context.Context) ([]models.Category, error) {
	cursor, err := c.categories.Find(ctx, bson.M{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	return castToSlice(ctx, cursor)
}

func (c *DatabaseRepository) Find(ctx context.Context, id string) (*models.Category, error) {
	var category models.Category

	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := c.decoder.Decode(c.categories.FindOne(ctx, bson.M{"_id": key}), &category); err != nil {
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
	defer cursor.Close(ctx)
	return castToSlice(ctx, cursor)
}

func (c *DatabaseRepository) Add(ctx context.Context, id, category string) error {
	_, err := c.categories.InsertOne(ctx, models.NewCategory(id, category))
	if err != nil {
		return err
	}
	return nil
}
