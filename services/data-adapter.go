package services

import (
	"github.com/gocraft/dbr"
	log "github.com/sirupsen/logrus"
)

type DataAdapter struct {
	connection *dbr.Connection
	tableName  string
}

func NewDataAdatper(conn *dbr.Connection, tableName string) *DataAdapter {
	return &DataAdapter{connection: conn}
}

func (r *DataAdapter) getSession() *dbr.Session {
	return r.connection.NewSession(nil)
}

func (r *DataAdapter) SetTableName(name string) {
	r.tableName = name
}

func (r *DataAdapter) GetTableName() string {
	return r.tableName
}

func (r *DataAdapter) Query(query string, value interface{}, result interface{}) error {
	session := r.getSession()

	stmt := session.Select("*").From(r.GetTableName()).Where(query, value).OrderBy("date")
	count, err := stmt.Load(result)
	if err != nil {
		log.WithFields(log.Fields{
			"query": stmt.Query,
			"error": err,
		}).Error("Error selecting records")
		return err
	} else {
		log.WithFields(log.Fields{
			"query": stmt.Query,
			"count": count,
		}).Debug("Selected records")
		return nil
	}
}

func (r *DataAdapter) QuerySingle(query string, value interface{}, result interface{}) error {
	session := r.getSession()

	stmt := session.Select("*").From(r.GetTableName()).Where(query, value).OrderBy("date")
	err := stmt.LoadOne(result)
	if err != nil {
		log.WithFields(log.Fields{
			"query": stmt.Query,
			"error": err,
		}).Error("Error selecting records")
		return err
	} else {
		log.WithFields(log.Fields{
			"query": stmt.Query,
		}).Debug("Retrived Record")
		return nil
	}
}

func (r *DataAdapter) Insert(columns []string, record interface{}) int64 {
	session := r.getSession()

	result := struct{ Id int64 }{}
	err := session.InsertInto(r.GetTableName()).Columns(columns...).Record(record).Returning("id").Load(&result)
	if err != nil {
		log.WithFields(log.Fields{
			"record": record,
			"error":  err,
		}).Panic("Error inserting record")
	} else {
		log.WithFields(log.Fields{
			"record": record,
		}).Debug("Selected Records")
	}
	if err != nil {
		log.WithFields(log.Fields{
			"record": record,
			"error":  err,
		}).Error("Error fetching created id")
	}
	return result.Id
}
