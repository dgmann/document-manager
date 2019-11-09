package filesystem

import (
	"context"
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

func (m *Manager) Index(ctx context.Context) (*Index, error) {
	if m.index != nil {
		return m.index, nil
	}

	filePath := filepath.Join(m.DataDirectory, "filesystem.gob")
	index, err := LoadIndexFromFile(filePath)
	if err != nil {
		m.index, err = CreateIndex(ctx, m.RecordDirectory)
		if err != nil {
			return nil, fmt.Errorf("error loading from filesystem: %w", err)
		}
	} else {
		m.index = index
	}
	defer func() {
		if err := m.index.Save(filePath); err != nil {
			logrus.WithError(err).Error("error saving filesystemindex to disk")
		}
	}()
	logrus.Info("load sub records")
	if err := m.index.LoadSubRecords(ctx, filepath.Join(m.DataDirectory, "splitted")); err != nil {
		return nil, err
	}

	return m.index, nil
}
