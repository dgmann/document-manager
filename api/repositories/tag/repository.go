package tag

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Repository interface {
	All() ([]string, error)
	ByPatient(id string) ([]string, error)
}

type DatabaseRepository struct {
	records *mgo.Collection
}

func NewDatabaseRepository(records *mgo.Collection) *DatabaseRepository {
	return &DatabaseRepository{records: records}
}

func (t *DatabaseRepository) All() ([]string, error) {
	var tags []string
	err := t.records.Find(nil).Distinct("tags", &tags)
	return tags, err
}

func (t *DatabaseRepository) ByPatient(id string) ([]string, error) {
	var tags []string
	err := t.records.Find(bson.M{ "patientId": id }).Distinct("tags", &tags)
	return tags, err
}
