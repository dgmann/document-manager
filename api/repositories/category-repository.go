package repositories

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"github.com/dgmann/document-manager/api/models"
)

type CategoryRepository struct {
	categories *mgo.Collection
	records    *mgo.Collection
}

func NewCategoryRepository(categories *mgo.Collection, records *mgo.Collection) *CategoryRepository {
	return &CategoryRepository{categories: categories, records: records}
}

func (c *CategoryRepository) All() ([]models.Category, error) {
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

func (c *CategoryRepository) Find(id string) (*models.Category, error) {
	var category models.Category

	if err := c.categories.FindId(bson.ObjectIdHex(id)).One(&category); err != nil {
		log.WithField("error", err).Panic("Cannot find record")
		return nil, err
	}
	return &category, nil
}

func (c *CategoryRepository) FindByPatient(id string) ([]models.Category, error) {
	categories := make([]models.Category, 0)
	var ids []bson.ObjectId
	if err := c.records.Find(bson.M{"patientId": id}).Distinct("categoryId", &ids); err != nil {
		log.WithField("error", err).Panic("Cannot find categories by patient id")
		return nil, err
	}
	if err := c.categories.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&categories); err != nil {
		log.WithField("error", err).Panic("Cannot resolve category ids")
		return nil, err
	}
	return categories, nil
}

func (c *CategoryRepository) Add(category string) error {
	return c.categories.Insert(models.NewCategory(category))
}
