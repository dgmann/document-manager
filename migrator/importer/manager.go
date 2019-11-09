package importer

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/api/client"
	"github.com/dgmann/document-manager/api/datastore"
	"github.com/dgmann/document-manager/migrator/categories"
	"github.com/dgmann/document-manager/migrator/records/filesystem"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"path"
	"path/filepath"
	"strconv"
)

type SubrecordIndex interface {
	SubRecords() filesystem.SubRecordList
}

type Manager struct {
	dataToImport        *Import
	dataDirectory       string
	importedRecordsPath string
	FilesystemManager   *filesystem.Manager
	importedRecords     map[string]ImportableRecord
	ImportErrors        ImportErrorList
	*Importer
	db *sqlx.DB
}

func NewManager(filesystemManager *filesystem.Manager, dataDirectory string, db *sqlx.DB, url string, retryCount int) *Manager {
	importer := NewImporter(url, retryCount)
	return &Manager{FilesystemManager: filesystemManager, dataDirectory: dataDirectory, Importer: importer, db: db, importedRecordsPath: filepath.Join(dataDirectory, "importedrecords.gob")}
}

func (m *Manager) DataToImport(ctx context.Context) (*Import, error) {
	if m.dataToImport != nil {
		return m.dataToImport, nil
	}

	if err := m.Load(FileName); err == nil {
		return m.dataToImport, nil
	}

	var filesToImport []ImportableRecord
	status := datastore.StatusDone
	index, err := m.FilesystemManager.Index(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching filesystem index: %w", err)
	}
	for _, r := range index.SubRecords() {
		subrecord := r.SubRecord()
		patId := strconv.Itoa(*subrecord.PatId)
		newRecord := client.NewRecord{
			ReceivedAt: subrecord.ReceivedAt,
			Date:       subrecord.Date,
			Status:     &status,
			PatientId:  &patId,
			Category:   subrecord.Spezialization,
		}
		filesToImport = append(filesToImport, ImportableRecord{
			NewRecord: newRecord,
			Path:      subrecord.Path,
		})
	}
	cats, err := categories.All(m.db)
	if err != nil {
		return nil, fmt.Errorf("error fetching categories: %w", err)
	}
	importData := Import{
		Records:    filesToImport,
		Categories: cats,
	}

	err = importData.Save(filepath.Join(m.dataDirectory, FileName))
	return &importData, err
}

func (m *Manager) Load(fileName string) error {
	importDataPath := filepath.Join(m.dataDirectory, fileName)
	var importData Import
	err := importData.Load(importDataPath)
	if err != nil {
		return err
	}
	m.dataToImport = &importData
	return nil
}

func (m *Manager) Import(ctx context.Context) error {
	dataToImport, err := m.DataToImport(ctx)
	if err != nil {
		return fmt.Errorf("error getting data to import: %w", err)
	}
	if err := m.importCategories(dataToImport.Categories); err != nil {
		return fmt.Errorf("error importing categories: %w", err)
	}
	if err := m.importRecords(ctx, dataToImport.Records); err != nil {
		return fmt.Errorf("error importing records: %w", err)
	}
	return nil
}

func (m *Manager) importCategories(cats []datastore.Category) error {
	for _, category := range cats {
		if err := m.Api.CreateCategory(category); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) ImportedRecords() map[string]ImportableRecord {
	if m.importedRecords == nil {
		importedRecords := make(map[string]ImportableRecord)
		_ = filesystem.LoadFromGob(importedRecords, m.importedRecordsPath)
		m.importedRecords = importedRecords
	}
	return m.importedRecords
}

func (m *Manager) importRecords(ctx context.Context, recordsToImport []ImportableRecord) (err error) {
	importedRecords := m.ImportedRecords()
	records := Difference(recordsToImport, importedRecords)
	importedCh, errCh := m.ImportRecords(ctx, records)

	defer func() {
		if err := filesystem.SaveToGob(importedRecords, m.importedRecordsPath); err != nil {
			logrus.Warn(err)
		}
	}()
	var importErrors ImportErrorList
	m.ImportErrors = importErrors
	defer func() {
		if len(importErrors) != 0 {
			recordsToReimport := make([]ImportableRecord, 0, len(importErrors))
			for i, r := range importErrors {
				recordsToImport[i] = *r.Record
			}
			reimportable := Import{Records: recordsToReimport}
			if err := reimportable.Save(path.Join(m.dataDirectory, "failedrecords.gob")); err != nil {
				logrus.WithError(err).Warn("error saving failedrecords.gob")
			}
			if err == nil {
				err = importErrors
			}
		}
	}()
	for {
		select {
		case record, ok := <-importedCh:
			if !ok {
				importedCh = nil
				break
			}
			importedRecords[record.Path] = *record
		case e, ok := <-errCh:
			if !ok {
				errCh = nil
				break
			}
			importErrors = append(importErrors, e)
		case <-ctx.Done():
			return nil
		}
		if importedCh == nil && errCh == nil {
			break
		}
	}
	return nil
}

func Difference(toImport []ImportableRecord, alreadyImported map[string]ImportableRecord) []ImportableRecord {
	var diff []ImportableRecord
	for _, record := range toImport {
		if _, ok := alreadyImported[record.Path]; !ok {
			diff = append(diff, record)
		}
	}
	return diff
}
