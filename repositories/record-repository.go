package repositories

import (
	"github.com/dgmann/document-manager-api/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

type RecordRepository struct {
	records *mgo.Collection
}

func NewRecordRepository(records *mgo.Collection) *RecordRepository {
	return &RecordRepository{records: records}
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

func (r *RecordRepository) Create(sender string) *models.Record {
	id := bson.NewObjectId()
	record := models.NewRecord(id, sender)
	if err := r.records.Insert(&record); err != nil {
		log.Panic(err)
	}
	return r.FindByObjectId(id)
}

func (r *RecordRepository) Delete(id string) error {
	key := bson.ObjectIdHex(id)
	err := r.records.RemoveId(key)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (r *RecordRepository) Update(id string, record models.Record) *models.Record {
	key := bson.ObjectIdHex(id)
	if err := r.records.UpdateId(key, bson.M{"$set": record}); err != nil {
		log.Panic(err)
	}
	return r.FindByObjectId(key)
}
