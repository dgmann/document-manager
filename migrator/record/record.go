package record

import (
	"strings"
	"path/filepath"
	"github.com/pkg/errors"
	"fmt"
	"strconv"
)

type Record struct {
	Name           string
	Path           string
	Spezialization string
	PatId          int
}

func (r *Record) String() string {
	return fmt.Sprintf("%s, %s: %s", r.PatId, r.Spezialization, r.Path)
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
	return &Record{
		Name:           fileName,
		Path:           path,
		Spezialization: spezialization,
		PatId:          patId,
	}, nil
}
