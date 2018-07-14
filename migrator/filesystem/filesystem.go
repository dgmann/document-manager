package filesystem

import (
	"path/filepath"
	"os"
	"github.com/pkg/errors"
	"github.com/dgmann/document-manager/migrator/shared"
	"io"
)

func CreateIndex(dir string) (*Index, error) {
	var files []CategorizableCloser
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, errors.Wrap(err, "directory to index does not exist")
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(info.Name()) != ".pdf" {
			return err
		}
		r, err := NewRecordFromPath(path)
		files = append(files, r)
		return err
	})
	index := newIndex(files)
	return index, err
}

type CategorizableCloser interface {
	shared.Categorizable
	io.Closer
}
