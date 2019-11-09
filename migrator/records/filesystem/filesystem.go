package filesystem

import (
	"context"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

func CreateIndex(ctx context.Context, dir string) (*Index, error) {
	var files []RecordContainerCloser
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, errors.Wrap(err, "directory to index does not exist")
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
			if info.IsDir() || filepath.Ext(info.Name()) != ".pdf" {
				return err
			}
			r, err := NewRecordFromPath(path)
			files = append(files, r)
			return err
		}
	})
	index := newIndex(files, dir)
	return index, err
}

func LoadIndexFromFile(dir string) (*Index, error) {
	index := &Index{}
	err := index.Load(dir)
	return index, err
}

type RecordContainerCloser interface {
	models.RecordContainer
	io.Closer
}
