package databasereader

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/dgmann/document-manager/migrator/shared"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	*shared.Manager
	index *Index
}

func NewManager(dbName, username, password, hostname, instance string) *Manager {
	return &Manager{shared.NewManager(dbName, username, password, hostname, instance), nil}
}

func (m *Manager) Index() (*Index, error) {
	if m.index != nil {
		return m.index, nil
	}

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
	m.index = newIndex(slice)
	return m.index, nil
}

func (m *Manager) loadRecords() ([]models.RecordContainer, error) {
	var records []*models.Record

	query := `select befund.Id, befund.Pat_Id, befund.Name, spez.Name as Category, befund.Pages from PdfDatabase.dbo.Befund as befund
			  JOIN PdfDatabase.dbo.Spezialisations as spez ON befund.Spezialisation_Id = spez.Id
 			  ORDER BY befund.Pat_Id`
	err := m.Db.Select(&records, query)
	if err != nil {
		return nil, errors.Wrap(err, "error loading records")
	}
	var recordContainers []models.RecordContainer
	for _, record := range records {
		recordContainers = append(recordContainers, record)
	}
	return recordContainers, nil
}

func (m *Manager) loadPdfs() ([]*PdfFile, error) {
	var pdfs []*PdfFile

	query := `select pdf.Id, pdf.Name, pdf.Timestamp, pdf.Date, pdf.State, pdf.SenderNr, pdf.Befund_Id, pdf.Pat_Id, pdf.Pages, type.Name as Type from PdfDatabase.dbo.PdfFiles as pdf
			  JOIN PdfDatabase.dbo.Type as type ON pdf.Type_Id = type.Id`
	err := m.Db.Select(&pdfs, query)
	if err != nil {
		return nil, errors.Wrap(err, "error loading pdfs")
	}
	return pdfs, nil
}

func mergeRecordsAndPdfs(records []models.RecordContainer, pdfs []*PdfFile) (map[int]models.RecordContainer, error) {
	idMap := createIdMap(records)
	var err error
	for _, pdf := range pdfs {
		if pdf.BefundId == nil {
			err = shared.WrapError(err, fmt.Sprintf("cannot assign pdf %d to any record. its record id is not set.", pdf.Id))
		} else if record, ok := idMap[*pdf.BefundId]; ok {
			r := record.Record()
			r.SubRecords = append(r.SubRecords, pdf.AsSubRecord())
			idMap[*pdf.BefundId] = r
		} else {
			err = shared.WrapError(err, fmt.Sprintf("cannot assign pdf %d to any record", pdf.Id))
		}
	}
	return idMap, err
}

func createIdMap(records []models.RecordContainer) map[int]models.RecordContainer {
	idMap := make(map[int]models.RecordContainer)
	for _, record := range records {
		idMap[record.Record().Id] = record
	}
	return idMap
}

func toSlice(records map[int]models.RecordContainer) []models.RecordContainer {
	v := make([]models.RecordContainer, 0, len(records))
	for _, value := range records {
		v = append(v, value)
	}
	return v
}
