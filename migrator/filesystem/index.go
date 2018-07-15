package filesystem

import (
	"github.com/dgmann/document-manager/migrator/shared"
	"io"
)

type Index struct {
	*shared.Index
}

func newIndex(data []CategorizableCloser) *Index {
	var cast []shared.Categorizable
	for _, d := range data {
		cast = append(cast, d)
	}
	return &Index{Index: shared.NewIndex("filesystem", cast)}
}

func (i *Index) Destroy() {
	for _, r := range i.GetAllCategorizable() {
		if closer, ok := r.(io.Closer); ok {
			closer.Close()
		}
	}
}

func (i *Index) GetRecords() []*shared.Record {
	var records []*shared.Record
	for _, r := range i.GetAllCategorizable() {
		record := r.(*Record)
		records = append(records, record.Record)
	}
	return records
}

func (i *Index) LoadAllSubRecords() error {
	var err error
	for _, r := range i.GetAllCategorizable() {
		record := r.(*Record)
		e := record.LoadSubRecords()
		if e != nil {
			err = shared.WrapError(err, e.Error())
		}
	}
	return err
}
