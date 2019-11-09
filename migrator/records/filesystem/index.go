package filesystem

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/migrator/records"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/dgmann/document-manager/migrator/splitter"
	"io"
	"path/filepath"
	"strings"
)

type Index struct {
	models.Index

	Path string
}

func newIndex(data []RecordContainerCloser, path string) *Index {
	var cast []models.RecordContainer
	for _, d := range data {
		cast = append(cast, d)
	}
	return &Index{Index: *models.NewIndex("filesystem", cast), Path: filepath.Clean(path)}
}

func (i *Index) Destroy() {
	for _, r := range i.Records() {
		if closer, ok := r.(io.Closer); ok {
			closer.Close()
		}
	}
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

func (i *Index) LoadSubRecords(ctx context.Context, dir string) error {
	recordCh, errCh := ToRecordChannel(ctx, i.Records())
	parallelCtx, cancel := context.WithCancel(ctx)
	go func() {
		for err := range errCh {
			if err != nil {
				cancel()
			}
		}
	}()

	err := records.Parallel(parallelCtx, recordCh, func(record models.RecordContainer) error {
		return loadSubRecord(record.(*Record), dir)
	})
	if len(err) > 0 {
		return errors.New(strings.Join(err, "; "))
	}
	return nil
}

func loadSubRecord(record *Record, dir string) error {
	if len(record.SubRecords) > 0 { // Already loaded
		return nil
	}
	data := bytes.NewReader(record.Content)
	subrecords, tmpDir, err := splitter.Split(data, dir)
	if err != nil {
		return fmt.Errorf("error splitting record. Path: %s, Error: %w", record.Path, err)
	}
	var convertedSubRecords []models.SubRecordContainer
	for _, subrecord := range subrecords {
		subrecord.BefundId = &record.Id
		subrecord.PatId = &record.PatId
		subrecord.Spezialization = &record.Spez
		convertedSubRecords = append(convertedSubRecords, &SubRecord{*subrecord})
	}
	record.SubRecords = convertedSubRecords
	record.SplittedPdfDir = tmpDir
	return nil
}

type SubRecordList []models.SubRecordContainer

func (i *Index) SubRecords() SubRecordList {
	var subrecords []models.SubRecordContainer
	for _, record := range i.Records() {
		subrecords = append(subrecords, record.Record().SubRecords...)
	}
	return subrecords
}

func (i *Index) Save(dir string) error {
	gob.Register(&Record{})
	gob.Register(&SubRecord{})
	return SaveToGob(i, dir)
}

func (i *Index) Load(dir string) error {
	gob.Register(&Record{})
	gob.Register(&SubRecord{})
	return LoadFromGob(i, dir)
}
