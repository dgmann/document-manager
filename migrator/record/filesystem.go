package record

import (
	"path"
	"strconv"
	"path/filepath"
	"os"
	"github.com/pkg/errors"
)

type Manager struct {
	directory string
}

func NewManager(dir string) *Manager {
	return &Manager{directory: dir}
}

func (m *Manager) CreateFileIndex() (*FileSystemIndex, error) {
	var files []*Record
	if _, err := os.Stat(m.directory); os.IsNotExist(err) {
		return nil, errors.Wrap(err, "directory to index does not exist")
	}

	err := filepath.Walk(m.directory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}
		r, err := NewRecordFromPath(path)
		files = append(files, r)
		return err
	})
	return NewFileSystemIndex(files), err
}

func (m *Manager) Load(patientId int, category string) (*RecordFile, error) {
	p := path.Join(m.directory, strconv.Itoa(patientId), category)
	return NewRecordFile(p)
}
