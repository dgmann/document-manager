package importer

import (
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api-client/record"
	"github.com/dgmann/document-manager/api-client/repository"
	"github.com/sirupsen/logrus"
	"os"
)

type Importer struct {
	repository       *repository.Repository
	recordRepository *record.Repository
	retryCount       int
}

func NewImporter(url string, retryCount int) *Importer {
	return &Importer{recordRepository: record.NewRepository(url), repository: repository.NewRepository(url), retryCount: retryCount}
}

func (i *Importer) ImportRecords(records []ImportableRecord) (<-chan *ImportableRecord, <-chan ImportError) {
	return parallel(records, i.uploadFunc())
}

func (i *Importer) Import(path string, model interface{}) error {
	return i.repository.Create(path, model)
}

func (i *Importer) uploadFunc() func(r *ImportableRecord) error {
	return func(r *ImportableRecord) error {
		f, err := os.Open(r.Path)
		if err != nil {
			return err
		}
		defer f.Close()

		logrus.WithField("record", r).Debug("upload file")
		for j := 0; j < i.retryCount; j++ {
			err = i.recordRepository.Create(&record.NewRecord{
				CreateRecord: r.CreateRecord,
				File:         f,
			})
			if err == nil {
				return nil
			}
			logrus.WithError(err).Errorf("error uploading file. Retry %d of %d", i, i.retryCount)
		}
		return errors.New(fmt.Sprintf("%s: %s", r.Path, err.Error()))
	}
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
