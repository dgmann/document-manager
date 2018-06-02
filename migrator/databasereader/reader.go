package databasereader

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	_ "github.com/denisenkom/go-mssqldb"
	"fmt"
	"github.com/dgmann/document-manager/migrator/shared"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	db *sqlx.DB
}

func (m *Manager) Open(dsn string) error {
	db, err := sqlx.Open("sqlserver", dsn)
	if err != nil {
		return err
	}
	m.db = db
	return nil
}

func (m *Manager) Close() {
	m.db.Close()
}

func (m *Manager) Load() (*Index, error) {
	records, err := m.loadRecords()
	if err != nil {
		return nil, err
	}
	pdfs, err := m.loadPdfs()
	if err != nil {
		return nil, err
	}
	merged, err := mergeRecordsAndPdfs(records, pdfs)
	if err != nil {
		log.WithField("error", err).Warning("There were errors assigning pdfs to records. This may happen for example when there are unprocessed pdfs in the inbox")
	}
	slice := toSlice(merged)
	index := newIndex(slice)
	return index, nil
}

func (m *Manager) loadRecords() ([]*shared.Record, error) {
	var records []*shared.Record

	query := `select befund.Id, befund.Pat_Id, spez.Name as Category from PdfDatabase.dbo.Befund as befund
			  JOIN PdfDatabase.dbo.Spezialisations as spez ON befund.Spezialisation_Id = spez.Id
 			  ORDER BY befund.Pat_Id`
	err := m.db.Select(&records, query)
	if err != nil {
		return nil, errors.Wrap(err, "error loading records")
	}
	return records, nil
}

func (m *Manager) loadPdfs() ([]*PdfFile, error) {
	var pdfs []*PdfFile

	query := `select pdf.Id, pdf.Name, pdf.Timestamp, pdf.Date, pdf.State, pdf.SenderNr, pdf.Befund_Id, pdf.Pat_Id, type.Name as Type from PdfDatabase.dbo.PdfFiles as pdf
			  JOIN PdfDatabase.dbo.Type as type ON pdf.Type_Id = type.Id`
	err := m.db.Select(&pdfs, query)
	if err != nil {
		return nil, errors.Wrap(err, "error loading pdfs")
	}
	return pdfs, nil
}

func mergeRecordsAndPdfs(records []*shared.Record, pdfs []*PdfFile) (map[int]*shared.Record, error) {
	idMap := createIdMap(records)
	var err error
	for _, pdf := range pdfs {
		if pdf.BefundId == nil {
			err = shared.WrapError(err, fmt.Sprintf("cannot assign pdf %d to any record. its record id is not set.", pdf.Id))
		} else if record, ok := idMap[*pdf.BefundId]; ok {
			record.SubRecords = append(record.SubRecords, pdf.AsSubRecord())
		} else {
			err = shared.WrapError(err, fmt.Sprintf("cannot assign pdf %d to any record", pdf.Id))
		}
	}
	return idMap, err
}

func createIdMap(records []*shared.Record) map[int]*shared.Record {
	idMap := make(map[int]*shared.Record)
	for _, record := range records {
		idMap[record.Id] = record
	}
	return idMap
}

func toSlice(records map[int]*shared.Record) []shared.Categorizable {
	v := make([]shared.Categorizable, 0, len(records))
	for _, value := range records {
		v = append(v, value)
	}
	return v
}
