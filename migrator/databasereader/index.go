package databasereader

import "github.com/dgmann/document-manager/migrator/shared"

type Index struct {
	*shared.Index
}

func (f *Index) GetRecords() []*shared.Record {
	var records []*shared.Record
	for _, r := range f.GetAllCategorizable() {
		c := r.(*shared.Record)
		records = append(records, c)
	}
	return records
}
