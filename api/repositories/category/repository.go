package category

import (
	"fmt"
	"github.com/dgmann/document-manager/api/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	All() ([]models.Category, error)
	Find(id string) (*models.Category, error)
	FindByPatient(id string) ([]models.Category, error)
	Add(id, category string) error
}

type DatabaseRepository struct {
	categories *mgo.Collection
	records    *mgo.Collection
}

func NewDatabaseRepository(categories, records *mgo.Collection) *DatabaseRepository {
	return &DatabaseRepository{categories: categories, records: records}
}

func (c *DatabaseRepository) All() ([]models.Category, error) {
	var categories []models.Category

	if err := c.categories.Find(bson.M{}).All(&categories); err != nil {
		log.Error(err)
		return nil, err
	}
	if categories == nil {
		categories = []models.Category{}
	}
	return categories, nil
}

func (c *DatabaseRepository) Find(id string) (*models.Category, error) {
	var category models.Category

	if err := c.categories.Find(bson.M{"_id": id}).One(&category); err != nil {
		log.WithField("error", err).Panic("Cannot find record")
		return nil, err
	}
	return &category, nil
}

func (c *DatabaseRepository) FindByPatient(id string) ([]models.Category, error) {
	categories := make([]models.Category, 0)
	var ids []string
	if err := c.records.Find(bson.M{"patientId": id}).Distinct("category", &ids); err != nil {
		log.WithField("error", err).Panic("Cannot find categories by patient id")
		return nil, err
	}
	fmt.Printf("ids %v", ids)
	if err := c.categories.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&categories); err != nil {
		log.WithField("error", err).Panic("Cannot resolve category ids")
		return nil, err
	}
	return categories, nil
}

func (c *DatabaseRepository) Add(id, category string) error {
	return c.categories.Insert(models.NewCategory(id, category))
}
