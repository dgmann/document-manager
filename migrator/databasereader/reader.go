package databasereader

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	_ "github.com/denisenkom/go-mssqldb"
	"fmt"
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

func (m *Manager) Load() (PatientCategoryRecordMap, error) {
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
		return nil, err
	}
	return groupByPatient(merged), nil
}

func (m *Manager) loadRecords() ([]*Record, error) {
	var records []*Record

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

	query := `select pdf.Id, pdf.Timestamp, pdf.Date, pdf.State, pdf.SenderNr, pdf.Befund_Id from PdfDatabase.dbo.PdfFiles as pdf`
	err := m.db.Select(&pdfs, query)
	if err != nil {
		return nil, errors.Wrap(err, "error loading pdfs")
	}
	return pdfs, nil
}

func mergeRecordsAndPdfs(records []*Record, pdfs []*PdfFile) (map[int]*Record, error) {
	idMap := createIdMap(records)
	var err error
	for _, pdf := range pdfs {
		if pdf.BefundId == nil {
			err = errors.Wrap(err, fmt.Sprintf("cannot assign pdf %d to any record. its record id is not set", pdf.Id))
		} else if record, ok := idMap[*pdf.BefundId]; ok {
			record.PdfFiles = append(record.PdfFiles, pdf)
		} else {
			err = errors.Wrap(err, fmt.Sprintf("cannot assign pdf %d to any record", pdf.Id))
		}
	}
	return idMap, err
}

func createIdMap(records []*Record) map[int]*Record {
	idMap := make(map[int]*Record)
	for _, record := range records {
		idMap[record.Id] = record
	}
	return idMap
}

func groupByPatient(records map[int]*Record) map[int]RecordsByCategory {
	patientMap := make(map[int]RecordsByCategory)
	for _, record := range records {
		if patientMap[record.PatId] == nil {
			patientMap[record.PatId] = make(RecordsByCategory)
		}
		patientMap[record.PatId][record.Category] = record
	}
	return patientMap
}
