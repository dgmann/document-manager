package importer

import (
	"os"
	"encoding/gob"
	"github.com/dgmann/document-manager/api-client/record"
	"github.com/dgmann/document-manager/migrator/patients"
	"github.com/dgmann/document-manager/migrator/categories"
)

type ImportableRecord struct {
	record.CreateRecord
	Path string
}

type Import struct {
	Categories []*categories.Category
	Patients   []*patients.Patient
	Records    []ImportableRecord
}

func (i Import) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(i)
}

func (i *Import) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(i)
}
