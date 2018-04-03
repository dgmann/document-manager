package repositories

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"github.com/dgmann/document-manager/api/models"
)

type CategoryRepository struct {
	categories *mgo.Collection
}

func NewCategoryRepository(categories *mgo.Collection) *CategoryRepository {
	return &CategoryRepository{categories: categories}
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

func (c *CategoryRepository) Add(category string) error {
	return c.categories.Insert(models.NewCategory(category))
}
