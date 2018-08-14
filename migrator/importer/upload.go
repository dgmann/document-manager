package importer

import (
	"github.com/dgmann/document-manager/api-client/record"
	"github.com/dgmann/document-manager/migrator/shared"
	"os"
	"github.com/sirupsen/logrus"
	"github.com/dgmann/document-manager/api-client/repository"
)

type Importer struct {
	repository       *repository.Repository
	recordRepository *record.Repository
}

func NewImporter(url string) *Importer {
	return &Importer{recordRepository: record.NewRepository(url), repository: repository.NewRepository(url)}
}

func (i *Importer) ImportRecords(records []ImportableRecord) []string {
	var interfaceSlice = make([]interface{}, len(records))
	for i, d := range records {
		interfaceSlice[i] = d
	}
	unsuccessFull := shared.Parallel(interfaceSlice, i.uploadFunc())
	return unsuccessFull
}

func (i *Importer) Import(path string, model interface{}) error {
	return i.repository.Create(path, model)
}

func (i *Importer) uploadFunc() shared.ParallelExecFunc {
	return func(value interface{}) error {
		r := value.(ImportableRecord)
		f, err := os.Open(r.Path)
		if err != nil {
			return err
		}
		logrus.WithField("record", r).Info("upload file")
		err = i.recordRepository.Create(&record.NewRecord{
			CreateRecord: r.CreateRecord,
			File:         f,
		})
		if err != nil {
			logrus.WithError(err).Error("error uploading file")
		}
		f.Close()
		return err
	}
}
