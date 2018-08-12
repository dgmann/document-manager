package importer

import (
	"github.com/dgmann/document-manager/api-client/record"
	"os"
	"encoding/gob"
)

type ImportableRecord struct {
	record.CreateRecord
	Path string
}

type ImportableRecordList []ImportableRecord

func (i ImportableRecordList) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(i)
}

func (i *ImportableRecordList) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(i)
}
