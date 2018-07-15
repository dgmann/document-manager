package filesystem

import (
	"strings"
	"path/filepath"
	"github.com/pkg/errors"
	"strconv"
	"github.com/dgmann/document-manager/migrator/splitter"
	"os"
	"github.com/dgmann/document-manager/migrator/records/models"
)

type embeddedRecord = models.Record

type Record struct {
	*embeddedRecord
	splittedPdfDir string
}

func NewRecordFromPath(dir string) (*Record, error) {
	directory, fileName := filepath.Split(dir)
	if directory == "" {
		return nil, errors.New("file dir cannot be parsed to create a record: " + dir)
	}

	spezialization := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	patIdString := filepath.Base(directory)
	patId, err := strconv.Atoi(patIdString)
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert patId to integer: "+dir)
	}
	r := &models.Record{
		Name:  fileName,
		Path:  dir,
		Spez:  spezialization,
		PatId: patId,
	}
	return &Record{embeddedRecord: r}, nil
}

func (r *Record) LoadSubRecords() error {
	subrecords, tmpDir, err := splitter.Split(r.Path)
	r.SubRecords = subrecords
	r.splittedPdfDir = tmpDir
	return err
}

func (r *Record) Close() error {
	return os.RemoveAll(r.splittedPdfDir)
}
