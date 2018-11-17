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
}

func NewImporter(url string) *Importer {
	return &Importer{recordRepository: record.NewRepository(url), repository: repository.NewRepository(url)}
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
		err = i.recordRepository.Create(&record.NewRecord{
			CreateRecord: r.CreateRecord,
			File:         f,
		})
		if err != nil {
			logrus.WithError(err).Error("error uploading file")
			return errors.New(fmt.Sprintf("%s: %s", r.Path, err.Error()))
		}
		return nil
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
