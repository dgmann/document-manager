package filesystem

import "github.com/dgmann/document-manager/migrator/records/models"

type EmbeddedSubrecord = models.SubRecord

type SubRecord struct {
	EmbeddedSubrecord
}

func (r *SubRecord) PageCount() int {
	count, err := getPageCount(r.Path)
	if err != nil {
		return -1
	}
	return count
}

func (r *SubRecord) SubRecord() *models.SubRecord {
	return &r.EmbeddedSubrecord
}
