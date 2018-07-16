package filesystem

import (
	"io"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/dgmann/document-manager/migrator/shared"
	"path/filepath"
)

type Index struct {
	*models.Index
	Path string
}

func newIndex(data []RecordContainerCloser, path string) *Index {
	var cast []models.RecordContainer
	for _, d := range data {
		cast = append(cast, d)
	}
	return &Index{Index: models.NewIndex("filesystem", cast), Path: filepath.Clean(path)}
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

func (i *Index) Validate() []string {
	invalidDirectories := make(map[string]struct{})
	for _, r := range i.Records() {
		dir := filepath.Dir(r.Record().Path)
		d := filepath.Dir(dir)
		if d != i.Path {
			invalidDirectories[dir] = struct{}{}
		}
	}
	keys := make([]string, 0, len(invalidDirectories))
	for k := range invalidDirectories {
		keys = append(keys, k)
	}
	return keys
}
