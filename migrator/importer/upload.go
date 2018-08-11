package importer

import (
	"github.com/dgmann/document-manager/api-client/record"
	"github.com/dgmann/document-manager/migrator/shared"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Importer struct {
	recordRepository *record.Repository
}

func NewImporter(url string) *Importer {
	return &Importer{recordRepository: record.NewRepository(url)}
}

func (i *Importer) Import(paths []string) []string {
	var interfaceSlice = make([]interface{}, len(paths))
	for i, d := range paths {
		interfaceSlice[i] = d
	}
	unsuccessFull := shared.Parallel(interfaceSlice, i.uploadFunc())
	return unsuccessFull
}

func (i *Importer) uploadFunc() shared.ParallelExecFunc {
	return func(value interface{}) error {
		p := value.(string)
		fileName := path.Base(p)
		receivedAt, err := time.Parse("02.01.2006", strings.TrimSuffix(fileName, filepath.Ext(fileName)))
		if err != nil {
			return err
		}
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		i.recordRepository.Create(&record.NewRecord{
			Sender:     "",
			File:       f,
			ReceivedAt: receivedAt,
		})
		f.Close()
		return nil
	}
}
