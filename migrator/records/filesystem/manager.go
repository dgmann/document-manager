package filesystem

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type Manager struct {
	RecordDirectory string
	DataDirectory   string
	index           *Index
}

func NewManager(recordDirectory, dataDirectory string) *Manager {
	return &Manager{RecordDirectory: recordDirectory, DataDirectory: dataDirectory}
}

func (m *Manager) Index() (*Index, error) {
	if m.index != nil {
		return m.index, nil
	}

	filePath := filepath.Join(m.DataDirectory, "filesystem.gob")
	index, err := LoadIndexFromFile(filePath)
	if err != nil {
		index, err = CreateIndex(m.RecordDirectory)
	}
	if err != nil {
		return nil, fmt.Errorf("error loading from filesystem: %w", err)
	}
	logrus.Info("load sub records")
	if err := index.LoadSubRecords(filepath.Join(m.DataDirectory, "splitted")); err != nil {
		return nil, err
	}
	m.index = index

	if err := m.index.Save(filePath); err != nil {
		logrus.WithError(err).Error("error saving filesystemindex to disk")
	}
	return m.index, nil
}
