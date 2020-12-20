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
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
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
	importer := NewImporter(url, retryCount, 1*time.Minute)
	return &Manager{FilesystemManager: filesystemManager, dataDirectory: dataDirectory, Importer: importer, db: db, importedRecordsPath: filepath.Join(dataDirectory, ImportedFileName)}
}

func (m *Manager) IsLoaded() bool {
	return m.dataToImport != nil
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

func (m *Manager) Files() []string {
	var files []string
	if _, err := os.Stat(filepath.Join(m.dataDirectory, FileName)); err == nil {
		files = append(files, FileName)
	}
	if _, err := os.Stat(filepath.Join(m.dataDirectory, FailedFileName)); err == nil {
		files = append(files, FailedFileName)
	}
	return files
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
		if err := filesystem.LoadFromGob(&importedRecords, m.importedRecordsPath); err != nil {
			logrus.Error(err)
		}
		m.importedRecords = importedRecords
	}
	return m.importedRecords
}

func (m *Manager) importRecords(ctx context.Context, recordsToImport []ImportableRecord) (err error) {
	importedRecords := m.ImportedRecords()
	records := Difference(recordsToImport, importedRecords)
	importedCh, errCh := m.ImportRecords(ctx, records)

	defer func() {
		for key, record := range importedRecords {
			record.File = nil
			importedRecords[key] = record
		}
		if err := filesystem.SaveToGob(importedRecords, m.importedRecordsPath); err != nil {
			logrus.Warn(err)
		}
	}()
	defer func() {
		if len(m.ImportErrors) != 0 {
			recordsToReimport := make([]ImportableRecord, 0, len(m.ImportErrors))
			for i, r := range m.ImportErrors {
				recordsToImport[i] = *r.Record
			}
			reimportable := Import{Records: recordsToReimport}
			if err := reimportable.Save(path.Join(m.dataDirectory, FailedFileName)); err != nil {
				logrus.WithError(err).Warn("error saving " + FailedFileName)
			}
			if err == nil {
				err = m.ImportErrors
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
			m.ImportErrors = append(m.ImportErrors, e)
			logrus.Errorf(e.Error())
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
	diff := make([]ImportableRecord, 0, len(toImport)-len(alreadyImported))
	for _, record := range toImport {
		if _, ok := alreadyImported[record.Path]; !ok {
			diff = append(diff, record)
		} else {
			logrus.Debugf("%s already imported. Skipping...", record.Path)
		}
	}
	return diff
}
