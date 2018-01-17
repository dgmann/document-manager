package services

import (
	"github.com/dgmann/document-manager-api/models"
	"github.com/gocraft/dbr"
	log "github.com/sirupsen/logrus"
)

const RECORDS_TABLE = "records"

type RecordService struct {
	connection *dbr.Connection
}

func NewRecordService(conn *dbr.Connection) *RecordService {
	service := RecordService{connection: conn}
	return &service
}

func (r *RecordService) getSession() *dbr.Session {
	return r.connection.NewSession(nil)
}

func (r *RecordService) Query(query string, value interface{}) []models.Record {
	session := r.getSession()
	defer session.Close()

	var records []models.Record
	stmt := session.Select("*").From(RECORDS_TABLE).Where("query", value).OrderBy("date")
	count, err := stmt.Load(records)
	if err != nil {
		log.WithFields(log.Fields{
			"query": stmt.Query,
			"error": err,
		}).Error("Error selecting records")
	} else {
		log.WithFields(log.Fields{
			"query": stmt.Query,
			"count": count,
		}).Debug("Selected Records")
	}
	return records
}

func (r *RecordService) QuerySingle(query string, value interface{}) models.Record {
	session := r.getSession()
	defer session.Close()

	var record models.Record
	stmt := session.Select("*").From(RECORDS_TABLE).Where("query", value).OrderBy("date")
	err := stmt.LoadOne(record)
	if err != nil {
		log.WithFields(log.Fields{
			"query": stmt.Query,
			"error": err,
		}).Error("Error selecting records")
	} else {
		log.WithFields(log.Fields{
			"query": stmt.Query,
		}).Debug("Retrived Record")
	}
	return record
}

func (r *RecordService) Insert(sender string, pages []models.Page) *models.Record {
	session := r.getSession()
	defer session.Close()

	insertColumns := []string{"sender", "pages"}
	record := &models.Record{
		Sender: sender,
		Pages:  pages,
	}
	stmt := session.InsertInto(RECORDS_TABLE).Columns(insertColumns...).Record(record)
	result, err := stmt.Exec()
	if err != nil {
		log.WithFields(log.Fields{
			"query":  stmt.Query,
			"record": record,
			"result": result,
			"error":  err,
		}).Error("Error inserting record")
	} else {
		log.WithFields(log.Fields{
			"query":  stmt.Query,
			"record": record,
			"result": result,
		}).Debug("Selected Records")
	}
	return record
}
