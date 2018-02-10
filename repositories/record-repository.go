package repositories

import (
	"github.com/dgmann/document-manager-api/models"
	"github.com/dgmann/document-manager-api/services"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"bytes"
)

type RecordRepository struct {
	records *mgo.Collection
	events  *services.EventService
	images  ImageRepository
}

func NewRecordRepository(records *mgo.Collection, images ImageRepository) *RecordRepository {
	return &RecordRepository{records: records, events: services.GetEventService(), images: images}
}

func (r *RecordRepository) Find(id string) *models.Record {
	return r.FindByObjectId(bson.ObjectIdHex(id))
}

func (r *RecordRepository) FindByObjectId(id bson.ObjectId) *models.Record {
	var record models.Record

	if err := r.records.FindId(id).One(&record); err != nil {
		log.WithField("error", err).Panic("Cannot find record")
	}
	record.Id = record.Primary.Hex()
	return &record
}

func (r *RecordRepository) Query(query interface{}) []*models.Record {
	var records []*models.Record

	if err := r.records.Find(query).All(&records); err != nil {
		log.Panic(err)
	}

	for _, record := range records {
		record.Id = record.Primary.Hex()
	}

	return records
}

func (r *RecordRepository) GetInbox() []*models.Record {
	return r.Query(bson.M{"processed": false})
}

func (r *RecordRepository) GetEscalated() []*models.Record {
	var records []*models.Record

	if err := r.records.Find(bson.M{"escalated": false}).All(&records); err != nil {
		log.Panic(err)
	}
	return records
}

func (r *RecordRepository) Create(sender string, images []*bytes.Buffer) (*models.Record, error) {
	id := bson.NewObjectId()

	record := models.NewRecord(id, sender)
	imageIds, err := r.images.Set(id.Hex(), images)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	record.SetPages(imageIds)

	if err := r.records.Insert(&record); err != nil {
		log.Error(err)
		r.images.Delete(id.Hex())
		return nil, err
	}
	created := r.FindByObjectId(id)

	r.events.Send(services.EventCreated, created)
	return created, nil
}

func (r *RecordRepository) Delete(id string) error {
	key := bson.ObjectIdHex(id)
	err := r.records.RemoveId(key)
	if err != nil {
		log.Error(err)
	}
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
