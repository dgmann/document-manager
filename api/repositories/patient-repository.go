package repositories

import (
	"github.com/globalsign/mgo"
	"github.com/dgmann/document-manager/api/models"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

type PatientRepository struct {
	patients *mgo.Collection
}

func NewPatientRepository(patients *mgo.Collection) *PatientRepository {
	return &PatientRepository{patients: patients}
}

func (p *PatientRepository) Add(patient *models.Patient) error {
	if _, err := p.patients.UpsertId(patient.Id, patient); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (c *PatientRepository) All() ([]*models.Patient, error) {
	patients := make([]*models.Patient, 0)

	if err := c.patients.Find(bson.M{}).All(&patients); err != nil {
		log.Error(err)
		return nil, err
	}

	return patients, nil
}

func (c *PatientRepository) FindByName(firstName, lastName string) ([]*models.Patient, error) {
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

func (c *PatientRepository) Find(id string) (*models.Patient, error) {
	var patient *models.Patient

	if err := c.patients.Find(bson.M{"_id": id}).One(&patient); err != nil {
		log.Error(err)
		return nil, err
	}

	return patient, nil
}
