package filesystem

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/sirupsen/logrus"
	pdf "github.com/unidoc/unidoc/pdf/model"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type EmbeddedRecord = models.Record

type Record struct {
	*EmbeddedRecord
	SplittedPdfDir string
	Content        []byte
}

func NewRecordFromPath(dir string) (*Record, error) {
	directory, fileName := filepath.Split(dir)
	if directory == "" {
		return nil, fmt.Errorf("file dir cannot be parsed to create a record: " + dir)
	}

	spezialization := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	patIdString := filepath.Base(directory)
	patId, err := strconv.Atoi(patIdString)
	if err != nil {
		return nil, fmt.Errorf("cannot convert patId to integer: %s. Error: %w", dir, err)
	}

	r := &models.Record{
		Name:  fileName,
		Path:  dir,
		Spez:  spezialization,
		PatId: patId,
		Pages: -2,
	}
	return &Record{EmbeddedRecord: r}, nil
}

func (r *Record) PageCount() int {
	if r.Pages != -2 {
		return r.Pages
	}
	pageCount, err := getPageCount(r.Path)
	if err != nil {
		logrus.Warn("cannot read page count of file: ", r.Path)
	}
	r.Pages = pageCount
	return r.Pages
}

func getPageCount(file string) (int, error) {
	f, err := os.Open(file)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return -1, err
	}

	return pdfReader.GetNumPages()
}

func (r *Record) GetPath() string {
	return r.Path
}

func (r *Record) Close() error {
	return os.RemoveAll(r.SplittedPdfDir)
}

func ToRecordChannel(ctx context.Context, values []models.RecordContainer) (chan models.RecordContainer, chan error) {
	workerCount := 4
	recordCh := make(chan models.RecordContainer, 10)
	errCh := make(chan error, workerCount)

	go func() {
		defer close(recordCh)
		defer close(errCh)
		for _, value := range values {
			select {
			case <-ctx.Done():
				return
			default:
				if record, ok := value.(*Record); ok {
					recordWithData := *record
					data, err := ioutil.ReadFile(recordWithData.Path)
					if err != nil {
						errCh <- err
						return
					}
					recordWithData.Content = data
					recordCh <- &recordWithData
				} else {
					panic(fmt.Errorf("error casting value to Record. %v", value))
				}
			}
		}
	}()

	return recordCh, errCh
}
