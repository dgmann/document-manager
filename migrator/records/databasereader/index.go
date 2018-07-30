package databasereader

import (
	"github.com/dgmann/document-manager/migrator/records/models"
	"encoding/gob"
)

type Index struct {
	models.Index
}

func newIndex(records []models.RecordContainer) *Index {
	return &Index{*models.NewIndex("database", records)}
}

func (i *Index) Save(dir string) error {
	gob.Register(&models.Record{})
	gob.Register(&PdfFile{})
	return i.Index.Save(dir)
}

func (i *Index) Load(path string) error {
	gob.Register(&models.Record{})
	gob.Register(&PdfFile{})
	return i.Index.Load(path)
}
