package category

import (
	"context"
	"github.com/dgmann/document-manager/api/models"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

func castToSlice(ctx context.Context, cursor mongo.Cursor) ([]models.Category, error) {
	categories := make([]models.Category, 0)

	for cursor.Next(ctx) {
		cat := models.Category{}
		if err := cursor.Decode(&cat); err != nil {
			logrus.WithError(err).Error("error decoding category from database")
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}
