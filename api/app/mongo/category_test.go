package mongo

import (
	"context"
	"encoding/hex"
	"github.com/dgmann/document-manager/api/app"
	"github.com/dgmann/document-manager/api/app/mock"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func buildCategoryService() (*CategoryService, *mock.Collection, *MockDecoder) {
	collection := mock.NewCollection()
	decoder := new(MockDecoder)
	repository := CategoryService{collection, collection, decoder}
	return &repository, collection, decoder
}

type MockDecoder struct {
	mock.Mock
}

func (d *MockDecoder) Decode(decodable Decodable, data interface{}) error {
	args := d.Called(decodable, data)
	return args.Error(0)
}

func TestDatabaseRepository_All(t *testing.T) {
	repository, collection, _ := buildCategoryService()

	ctx := context.Background()
	cursor := mock.NewCursor()

	collection.On("Find", ctx, primitive.M{}).Once().Return(cursor, nil)
	cursor.On("Next", ctx).Twice().Return(true)
	cursor.On("Next", ctx).Once().Return(false)
	cursor.On("Decode", mock.Anything).Return(nil)
	cursor.On("Close", ctx).Return(nil)

	res, err := repository.All(ctx)

	collection.AssertExpectations(t)
	assert.Nil(t, err, "no error expected")
	assert.Len(t, res, 2, "wrong length of result")
}

func TestDatabaseRepository_Find(t *testing.T) {
	repository, collection, decoder := buildCategoryService()
	id := "4ecc05e55dd98a436ddcc47c"
	objectID, _ := primitive.ObjectIDFromHex(id)

	ctx := context.Background()
	result := &mongo.SingleResult{}

	collection.On("FindOne", ctx, primitive.M{"_id": objectID}).Once().Return(result)
	decoder.On("Decode", result, mock.Anything).Return(nil)

	cat, err := repository.Find(ctx, id)

	collection.AssertExpectations(t)
	assert.Nil(t, err, "no error expected")
	assert.NotNil(t, cat, "result is nil")
}

func TestDatabaseRepository_Find_InvalidId(t *testing.T) {
	repository, collection, _ := buildCategoryService()
	id := "1"

	ctx := context.Background()

	_, err := repository.Find(ctx, id)

	collection.AssertExpectations(t)
	assert.Equal(t, err, hex.ErrLength)
}

func TestDatabaseRepository_FindByPatient(t *testing.T) {
	repository, collection, _ := buildCategoryService()
	patientID := "1"
	categoryNames := []interface{}{"cat1", "cat2"}

	ctx := context.Background()
	cursor := mock.NewCursor()

	collection.On("Distinct", ctx, "category", primitive.M{"patientId": patientID}).Once().Return(categoryNames, nil)
	collection.On("Find", ctx, primitive.M{"_id": primitive.M{"$in": categoryNames}}).Once().Return(cursor, nil)
	cursor.On("Next", ctx).Twice().Return(true)
	cursor.On("Next", ctx).Once().Return(false)
	cursor.On("Decode", mock.Anything).Return(nil)
	cursor.On("Close", ctx).Return(nil)

	cats, err := repository.FindByPatient(ctx, patientID)

	collection.AssertExpectations(t)
	assert.Nil(t, err, "no error expected")
	assert.Len(t, cats, 2, "wrong length")
}

func TestDatabaseRepository_Add(t *testing.T) {
	repository, collection, _ := buildCategoryService()
	id := "cat1"
	name := "Category 1"

	ctx := context.Background()

	collection.On("InsertOne", ctx, app.NewCategory(id, name)).Once().Return(&mongo.InsertOneResult{InsertedID: id}, nil)

	err := repository.Add(ctx, id, name)

	collection.AssertExpectations(t)
	assert.Nil(t, err, "no error expected")
}
