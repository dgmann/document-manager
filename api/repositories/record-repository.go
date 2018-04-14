package repositories

import (
	"github.com/dgmann/document-manager/api/models"
	"github.com/dgmann/document-manager/api/services"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"github.com/dgmann/document-manager/shared"
)

type RecordRepository struct {
	records *mgo.Collection
	events  *services.EventService
	images  ImageRepository
}

func NewRecordRepository(records *mgo.Collection, images ImageRepository) *RecordRepository {
	processedIndex := mgo.Index{
		Key:        []string{"patientId", "-date", "tags"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}

	err := records.EnsureIndex(processedIndex)
	if err != nil {
		log.Panicf("Error setting indices %s", err)
	}

	actionIndex := mgo.Index{
		Key:        []string{"requiredAction"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}

	err = records.EnsureIndex(actionIndex)
	if err != nil {
		log.Panicf("Error setting indices %s", err)
	}
	return &RecordRepository{records: records, events: services.GetEventService(), images: images}
}

func (r *RecordRepository) All() ([]*models.Record, error) {
	return r.Query(bson.M{})
}

func (r *RecordRepository) Find(id string) *models.Record {
	return r.FindByObjectId(bson.ObjectIdHex(id))
}

func (r *RecordRepository) FindByObjectId(id bson.ObjectId) *models.Record {
	var record models.Record

	if err := r.records.FindId(id).One(&record); err != nil {
		log.WithField("error", err).Panic("Cannot find record")
	}
	return &record
}

func (r *RecordRepository) FindByPatientId(id string) ([]*models.Record, error) {
	records, err := r.Query(bson.M{"patientId": id})
	if err != nil {
		log.WithField("error", err).Panic("Cannot find records by patient id")
		return nil, err
	}
	return records, nil
}

func (r *RecordRepository) Query(query map[string]interface{}) ([]*models.Record, error) {
	records := make([]*models.Record, 0)

	if err := r.records.Find(query).All(&records); err != nil {
		log.Error(err)
		return nil, err
	}

	return records, nil
}

func (r *RecordRepository) GetInbox() ([]*models.Record, error) {
	return r.Query(bson.M{"$or": []bson.M{{"date": nil}, {"patientId": ""}, {"categoryId": nil}}})
}

func (r *RecordRepository) GetEscalated() []*models.Record {
	var records []*models.Record

	if err := r.records.Find(bson.M{"escalated": false}).All(&records); err != nil {
		log.Panic(err)
	}
	return records
}

func (r *RecordRepository) Create(sender string, images []*shared.Image) (*models.Record, error) {
	record := models.NewRecord(sender)
	imageIds, err := r.images.Set(record.Id.Hex(), images)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	record.SetPages(imageIds)

	if err := r.records.Insert(&record); err != nil {
		log.Error(err)
		r.images.Delete(record.Id.Hex())
		return nil, err
	}
	created := r.FindByObjectId(record.Id)

	r.events.Send(services.EventCreated, created)
	return created, nil
}

func (r *RecordRepository) Delete(id string) error {
	key := bson.ObjectIdHex(id)
	err := r.records.RemoveId(key)
	if err != nil {
		log.Error(err)
	}
	r.images.Delete(id)
	r.events.Send(services.EventDeleted, id)
	return err
}

func (r *RecordRepository) Update(id string, record models.Record) *models.Record {
	key := bson.ObjectIdHex(id)
	if err := r.records.UpdateId(key, bson.M{"$set": record}); err != nil {
		log.Panic(err)
	}
	updated := r.FindByObjectId(key)
	r.events.Send(services.EventUpdated, updated)
	return updated
}
