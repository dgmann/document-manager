package importer

import (
	client "github.com/dgmann/document-manager/apiclient"
	"github.com/dgmann/document-manager/migrator/records/filesystem"
	"github.com/dgmann/document-manager/migrator/records/models"
	"strconv"
)

const FileName string = "datatoimport.gob"
const FailedFileName string = "failedrecords.gob"
const ImportedFileName string = "importedrecords.gob"

type ImportableRecord struct {
	client.NewRecord
	Path string
}

func (i *ImportableRecord) Record() *models.Record {
	return &models.Record{
		Id:         -1,
		Name:       "",
		PatId:      i.PatientId(),
		Spez:       *i.Category,
		Pages:      i.PageCount(),
		Path:       i.Path,
		SubRecords: nil,
	}
}

func (i *ImportableRecord) Spezialization() string {
	return *i.Category
}

func (i *ImportableRecord) PatientId() int {
	if res, err := strconv.Atoi(*i.NewRecord.PatientId); err != nil {
		return res
	} else {
		return 0
	}
}

func (i *ImportableRecord) PageCount() int {
	return 0
}

type Import struct {
	Categories []client.Category
	Records    []ImportableRecord
}

func (i *Import) Save(path string) error {
	return filesystem.SaveToGob(i, path)
}

func (i *Import) Load(path string) error {
	return filesystem.LoadFromGob(i, path)
}
