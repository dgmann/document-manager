package filesystem

import (
	"io"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/dgmann/document-manager/migrator/shared"
)

type Index struct {
	*models.Index
}

func newIndex(data []RecordContainerCloser) *Index {
	var cast []models.RecordContainer
	for _, d := range data {
		cast = append(cast, d)
	}
	return &Index{Index: models.NewIndex("filesystem", cast)}
}

func (i *Index) Destroy() {
	for _, r := range i.Records() {
		if closer, ok := r.(io.Closer); ok {
			closer.Close()
		}
	}
}

func (i *Index) LoadAllSubRecords() error {
	var err error
	for _, r := range i.Records() {
		e := r.Record().LoadSubRecords()
		if e != nil {
			err = shared.WrapError(err, e.Error())
		}
	}
	return err
}
