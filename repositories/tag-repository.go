package repositories

import "github.com/globalsign/mgo"

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
