package services

import (
	"github.com/gocraft/dbr"
)

const RECORDS_TABLE = "records"

func NewRecordDataAdapter(conn *dbr.Connection) *DataAdapter {
	return &DataAdapter{connection: conn, tableName: RECORDS_TABLE}
}
