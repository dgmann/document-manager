package watcher

import (
	"github.com/dgmann/document-manager/directory-watcher/parser"
	"io/ioutil"
	"os"
	"time"
	log "github.com/sirupsen/logrus"
	"errors"
	"fmt"
	"path"
	"github.com/dgmann/document-manager/api-client/record"
)

type DirectoryWatcher struct {
	ticker        *time.Ticker
	watchedFiles  map[string]struct{}
	recordChannel chan *record.NewRecord
	retryCount    int
}

func NewDirectoryWatcher(scanInterval, retry int) *DirectoryWatcher {
	return &DirectoryWatcher{
		ticker:        time.NewTicker(time.Duration(scanInterval) * time.Second),
		recordChannel: make(chan *record.NewRecord),
		watchedFiles:  make(map[string]struct{}),
		retryCount:    retry,
	}
}

func (w *DirectoryWatcher) Close() {
	w.ticker.Stop()
	close(w.recordChannel)
}

func (w *DirectoryWatcher) Watch(dir string, parser parser.Parser) <-chan *record.NewRecord {
	go func() {
		for range w.ticker.C {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				if f.IsDir() || path.Ext(f.Name()) != ".pdf" {
					continue
				}

				record := parser.Parse(f.Name())
				record.PdfPath = path.Join(dir, f.Name())
				w.add(record)
			}
		}
	}()
	return w.recordChannel
}

func (w *DirectoryWatcher) Done(record *record.NewRecord) {
	if err := w.remove(record); err != nil {
		log.WithField("error", err).WithField("record", record).Infof("error processing record")
	} else {
		log.WithField("record", record).Infof("record sucessfully processed")
	}
}

func (w *DirectoryWatcher) Error(record *record.NewRecord) {
	record.RetryCounter++
	if record.RetryCounter <= w.retryCount {
		go func(record *models.RecordCreate) {
			time.Sleep(2 * time.Second)
			w.recordChannel <- record
			log.WithField("record", record).Info("requeue record")
		}(record)
	} else {
		log.WithField("record", record).WithField("retryCount", w.retryCount).Info("retry counter exceeded")
	}
}

func (w *DirectoryWatcher) add(record *record.NewRecord) {
	if _, ok := w.watchedFiles[record.PdfPath]; ok {
		return
	}

	w.watchedFiles[record.PdfPath] = struct{}{}

	w.recordChannel <- record
}

func (w *DirectoryWatcher) remove(record *record.NewRecord) error {
	if _, ok := w.watchedFiles[record.PdfPath]; !ok {
		return errors.New(fmt.Sprintf("record %s not found", record.PdfPath))
	}

	if err := os.Remove(record.PdfPath); err != nil {
		return err
	}

	delete(w.watchedFiles, record.PdfPath)
	return nil
}
