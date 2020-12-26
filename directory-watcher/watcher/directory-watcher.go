package watcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	client "github.com/dgmann/document-manager/apiclient"
	"github.com/dgmann/document-manager/directory-watcher/parser"
	log "github.com/sirupsen/logrus"
)

type NewRecord struct {
	*client.NewRecord
	PdfPath      string
	RetryCounter int
}

func NewNewRecord(record *client.NewRecord, pdfPath string) *NewRecord {
	return &NewRecord{NewRecord: record, PdfPath: pdfPath, RetryCounter: 0}
}

type DirectoryWatcher struct {
	ticker        *time.Ticker
	watchedFiles  map[string]struct{}
	recordChannel chan *NewRecord
	retryCount    int
}

func NewDirectoryWatcher(scanInterval, retry int) *DirectoryWatcher {
	return &DirectoryWatcher{
		ticker:        time.NewTicker(time.Duration(scanInterval) * time.Second),
		recordChannel: make(chan *NewRecord),
		watchedFiles:  make(map[string]struct{}),
		retryCount:    retry,
	}
}

func (w *DirectoryWatcher) Close() {
	w.ticker.Stop()
	close(w.recordChannel)
}

func (w *DirectoryWatcher) Watch(dir string, parser parser.Parser) <-chan *NewRecord {
	go func() {
		for range w.ticker.C {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Errorf("error reading directory %s: %s", dir, err)
				continue
			}

			for _, f := range files {
				if f.IsDir() || path.Ext(f.Name()) != ".pdf" {
					continue
				}

				parsed := parser.Parse(f.Name())
				record := NewNewRecord(parsed, f.Name())
				record.PdfPath = path.Join(dir, f.Name())
				w.add(record)
			}
		}
	}()
	return w.recordChannel
}

func (w *DirectoryWatcher) Done(record *NewRecord) {
	if err := w.remove(record); err != nil {
		log.WithField("error", err).WithField("path", record.PdfPath).Infof("could not remove processed record")
	} else {
		log.WithField("path", record.PdfPath).Infof("record sucessfully processed")
	}
}

func (w *DirectoryWatcher) Error(record *NewRecord) {
	record.RetryCounter++
	if record.RetryCounter <= w.retryCount {
		go func(record *NewRecord) {
			time.Sleep(2 * time.Second)
			w.recordChannel <- record
			log.WithField("path", record.PdfPath).Info("requeue record")
		}(record)
	} else {
		log.WithField("path", record.PdfPath).WithField("retryCount", w.retryCount).Info("retry counter exceeded")
	}
}

func (w *DirectoryWatcher) add(record *NewRecord) {
	if _, ok := w.watchedFiles[record.PdfPath]; ok {
		return
	}

	w.watchedFiles[record.PdfPath] = struct{}{}

	w.recordChannel <- record
}

func (w *DirectoryWatcher) remove(record *NewRecord) error {
	if _, ok := w.watchedFiles[record.PdfPath]; !ok {
		return errors.New(fmt.Sprintf("record %s not found", record.PdfPath))
	}

	if err := os.Remove(record.PdfPath); err != nil {
		return err
	}

	delete(w.watchedFiles, record.PdfPath)
	return nil
}
