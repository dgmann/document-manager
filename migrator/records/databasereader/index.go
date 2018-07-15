package databasereader

import "github.com/dgmann/document-manager/migrator/records/models"

type Index struct {
	*models.Index
}

func newIndex(records []models.RecordContainer) *Index {
	return &Index{models.NewIndex("database", records)}
}
