package patient

import (
	"github.com/dgmann/document-manager/api/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	Add(patient *models.Patient) error
	All() ([]*models.Patient, error)
	FindByName(firstName, lastName string) ([]*models.Patient, error)
	Find(id string) (*models.Patient, error)
}

type DatabaseRepository struct {
	patients *mgo.Collection
}

func NewDatabaseRepository(patients *mgo.Collection) *DatabaseRepository {
	return &DatabaseRepository{patients: patients}
}

func (p *DatabaseRepository) Add(patient *models.Patient) error {
	if _, err := p.patients.UpsertId(patient.Id, patient); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (c *DatabaseRepository) All() ([]*models.Patient, error) {
	patients := make([]*models.Patient, 0)

	if err := c.patients.Find(bson.M{}).All(&patients); err != nil {
		log.Error(err)
		return nil, err
	}

	return patients, nil
}

func (c *DatabaseRepository) FindByName(firstName, lastName string) ([]*models.Patient, error) {
	patients := make([]*models.Patient, 0)

	if err := c.patients.Find(bson.M{
		"firstName": &bson.RegEx{Pattern: firstName, Options: "i"},
		"lastName":  &bson.RegEx{Pattern: lastName, Options: "i"},
	}).All(&patients); err != nil {
		log.Error(err)
		return nil, err
	}

	return patients, nil
}

func (c *DatabaseRepository) Find(id string) (*models.Patient, error) {
	var patient *models.Patient

	if err := c.patients.Find(bson.M{"_id": id}).One(&patient); err != nil {
		log.Error(err)
		return nil, err
	}

	return patient, nil
}
