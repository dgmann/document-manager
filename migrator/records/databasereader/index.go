package databasereader

import (
	"encoding/gob"
	"github.com/dgmann/document-manager/migrator/records/filesystem"
	"github.com/dgmann/document-manager/migrator/records/models"
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
	return filesystem.SaveToGob(i, dir)
}

func (i *Index) Load(path string) error {
	gob.Register(&models.Record{})
	gob.Register(&PdfFile{})
	return filesystem.LoadFromGob(i, path)
}
