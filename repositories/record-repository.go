package repositories

import (
	"github.com/dgmann/document-manager-api/models"
	"github.com/dgmann/document-manager-api/services"
)

type RecordRepository struct {
	records *services.RecordService
}

func NewRecordRepository(service *services.RecordService) *RecordRepository {
	return &RecordRepository{records: service}
}

func (r *RecordRepository) Find(id string) models.Record {
	record := r.records.QuerySingle("id = ?", id)
	return record
}

func (r *RecordRepository) GetInbox() []models.Record {
	records := r.records.Query("processed = ?", false)
	return records
}

func (r *RecordRepository) GetEscalated() []models.Record {
	records := r.records.Query("escalated = ?", true)
	return records
}

func (r *RecordRepository) Create(sender string) *models.Record {
	records := r.records.Insert(sender, nil)
	return records
}
