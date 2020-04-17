package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type TagService struct {
	records distinctFinder
}

func NewTagService(records distinctFinder) *TagService {
	return &TagService{records: records}
}

func (t *TagService) All(ctx context.Context) ([]string, error) {
	return t.query(ctx, bson.M{})
}

func (t *TagService) ByPatient(ctx context.Context, id string) ([]string, error) {
	return t.query(ctx, bson.M{"patientId": id})
}

func (t *TagService) query(ctx context.Context, query interface{}) ([]string, error) {
	res, err := t.records.Distinct(ctx, "tags", bson.M{})
	if err != nil {
		return []string{}, err
	}
	tags := make([]string, 0, len(res))
	for _, tag := range res {
		tags = append(tags, tag.(string))
	}
	return tags, err
}
