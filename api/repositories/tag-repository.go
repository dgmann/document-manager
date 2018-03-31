package repositories

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type TagRepository struct {
	records *mgo.Collection
}

func NewTagRepository(records *mgo.Collection) *TagRepository {
	return &TagRepository{records:records}
}

func (t* TagRepository) All() ([]string, error) {
	var tags []string
	err := t.records.Find(nil).Distinct("tags", &tags)
	return tags, err
}

func (t* TagRepository) ByPatient(id string) ([]string, error) {
	var tags []string
	err := t.records.Find(bson.M{ "patientId": id }).Distinct("tags", &tags)
	return tags, err
}
