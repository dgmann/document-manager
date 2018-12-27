package mongo

import (
	"context"
	"github.com/dgmann/document-manager/api/app"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

type CategoryService struct {
	categories categoryCollection
	records    distinctFinder
	decoder    Decoder
}

type categoryCollection interface {
	finder
	oneFinder
	oneInserter
}

func NewCategoryService(categories categoryCollection, records distinctFinder) *CategoryService {
	return &CategoryService{categories: categories, records: records, decoder: NewDefaultDecoder()}
}

func (c *CategoryService) All(ctx context.Context) ([]app.Category, error) {
	cursor, err := c.categories.Find(ctx, bson.M{})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	return castToCategorySlice(ctx, cursor)
}

func (c *CategoryService) Find(ctx context.Context, id string) (*app.Category, error) {
	var category app.Category

	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := c.decoder.Decode(c.categories.FindOne(ctx, bson.M{"_id": key}), &category); err != nil {
		logrus.WithField("error", err).Panic("Cannot find record")
		return nil, err
	}
	return &category, nil
}

func (c *CategoryService) FindByPatient(ctx context.Context, id string) ([]app.Category, error) {
	ids, err := c.records.Distinct(ctx, "category", bson.M{"patientId": id})
	if err != nil {
		logrus.WithField("error", err).Panic("Cannot find categories by patient id")
		return nil, err
	}

	cursor, err := c.categories.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		logrus.WithField("error", err).Panic("Cannot resolve category ids")
		return nil, err
	}
	defer cursor.Close(ctx)
	return castToCategorySlice(ctx, cursor)
}

func (c *CategoryService) Add(ctx context.Context, id, category string) error {
	_, err := c.categories.InsertOne(ctx, app.NewCategory(id, category))
	if err != nil {
		return err
	}
	return nil
}

func castToCategorySlice(ctx context.Context, cursor mongo.Cursor) ([]app.Category, error) {
	categories := make([]app.Category, 0)

	for cursor.Next(ctx) {
		cat := app.Category{}
		if err := cursor.Decode(&cat); err != nil {
			logrus.WithError(err).Error("error decoding category from database")
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}
