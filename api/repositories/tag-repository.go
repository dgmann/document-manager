package repositories

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type TagRepository interface {
	All() ([]string, error)
	ByPatient(id string) ([]string, error)
}

type DBTagRepository struct {
	records *mgo.Collection
}

func newDBTagRepository(records *mgo.Collection) *DBTagRepository {
	return &DBTagRepository{records: records}
}

func (t *DBTagRepository) All() ([]string, error) {
	var tags []string
	err := t.records.Find(nil).Distinct("tags", &tags)
	return tags, err
}

func (t *DBTagRepository) ByPatient(id string) ([]string, error) {
	var tags []string
	err := t.records.Find(bson.M{ "patientId": id }).Distinct("tags", &tags)
	return tags, err
}
