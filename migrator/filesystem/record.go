package filesystem

import (
	"strings"
	"path/filepath"
	"github.com/pkg/errors"
	"strconv"
	"github.com/dgmann/document-manager/migrator/shared"
	"github.com/dgmann/document-manager/migrator/splitter"
	"os"
)

type Record struct {
	*shared.Record
	splittedPdfDir string
}

func NewRecordFromPath(path string) (*Record, error) {
	split := strings.Split(path, "/")
	if len(split) < 2 {
		return nil, errors.New("file path cannot be parsed to create a record: " + path)
	}

	fileName := split[len(split)-1]
	spezialization := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	patIdString := split[len(split)-2]
	patId, err := strconv.Atoi(patIdString)
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert patId to integer: "+path)
	}
	r := &shared.Record{
		Name:           fileName,
		Path:           path,
		Spezialization: spezialization,
		PatId:          patId,
	}
	return &Record{Record: r}, nil
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
