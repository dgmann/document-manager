package repositories

import (
	"github.com/dgmann/document-manager-api/models"
	"github.com/dgmann/document-manager-api/services"
	log "github.com/sirupsen/logrus"
)

type RecordRepository struct {
	data *services.DataAdapter
}

func NewRecordRepository(data *services.DataAdapter) *RecordRepository {
	return &RecordRepository{data: data}
}

func (r *RecordRepository) Find(id int64) *models.Record {
	var record models.Record
	if r.data.QuerySingle("id = ?", id, &record) != nil {
		log.Panic("Cannot find record")
	}
	return &record
}

func (r *RecordRepository) GetInbox() []*models.Record {
	var records []*models.Record
	if err := r.data.Query("processed = ?", false, &records); err != nil {
		log.Panic(err)
	}
	return records
}

func (r *RecordRepository) GetEscalated() []*models.Record {
	var records []*models.Record

	if r.data.Query("escalated = ?", true, &records) != nil {
		log.Panic("Cannot find escalated records")
	}
	return records
}

func (r *RecordRepository) Create(sender string) *models.Record {
	record := models.Record{Sender: sender}
	id := r.data.Insert([]string{"sender"}, record)
	return r.Find(id)
}
