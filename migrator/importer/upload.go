package importer

import (
	"github.com/dgmann/document-manager/api-client/record"
	"github.com/dgmann/document-manager/migrator/shared"
	"os"
)

type Importer struct {
	recordRepository *record.Repository
}

func NewImporter(url string) *Importer {
	return &Importer{recordRepository: record.NewRepository(url)}
}

func (i *Importer) Import(records ImportableRecordList) []string {
	var interfaceSlice = make([]interface{}, len(records))
	for i, d := range records {
		interfaceSlice[i] = d
	}
	unsuccessFull := shared.Parallel(interfaceSlice, i.uploadFunc())
	return unsuccessFull
}

func (i *Importer) uploadFunc() shared.ParallelExecFunc {
	return func(value interface{}) error {
		r := value.(ImportableRecord)
		f, err := os.Open(r.Path)
		if err != nil {
			return err
		}
		i.recordRepository.Create(&record.NewRecord{
			CreateRecord: r.CreateRecord,
			File:         f,
		})
		f.Close()
		return nil
	}
}
