package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
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
		return nil, fmt.Errorf("while loading categories: %w", err)
	}
	err = cursor.Close(ctx)
	if err != nil {
		return nil, fmt.Errorf("while closing cursor: %w", err)
	}
	return castToCategorySlice(ctx, cursor)
}

func (c *CategoryService) Find(ctx context.Context, id string) (*app.Category, error) {
	var category app.Category

	key, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, app.NewNotFoundError(id, Categories, err)
	}
	res := c.categories.FindOne(ctx, bson.M{"_id": key})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := c.decoder.Decode(res, &category); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, app.NewNotFoundError(id, Categories, err)
		}
		return nil, err
	}
	return &category, nil
}

func (c *CategoryService) FindByPatient(ctx context.Context, id string) ([]app.Category, error) {
	ids, err := c.records.Distinct(ctx, "category", bson.M{"patientId": id})
	if err != nil {
		return nil, fmt.Errorf("while finding categories by patient id: %w", err)
	}

	cursor, err := c.categories.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, fmt.Errorf("while resolving categories %v: %w", ids, err)
	}
	err = cursor.Close(ctx)
	if err != nil {
		return nil, fmt.Errorf("while closing cursor: %w", err)
	}
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
			return nil, fmt.Errorf("decoding cartegories from database: %w", err)
		}
		categories = append(categories, cat)
	}

	return categories, nil
}
