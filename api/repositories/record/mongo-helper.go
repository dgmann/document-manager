package record

import (
	"context"
	"github.com/dgmann/document-manager/api/models"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

func castToSlice(ctx context.Context, cursor mongo.Cursor) ([]*models.Record, error) {
	records := make([]*models.Record, 0)

	for cursor.Next(ctx) {
		r := models.Record{}
		if err := cursor.Decode(&r); err != nil {
			logrus.WithError(err).Error("error decoding category from database")
			return nil, err
		}
		records = append(records, &r)
	}

	return records, nil
}
