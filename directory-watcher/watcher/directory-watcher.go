package watcher

import (
	"github.com/dgmann/document-manager/directory-watcher/models"
	"github.com/dgmann/document-manager/directory-watcher/parser"
	"io/ioutil"
	"os"
	"time"
	log "github.com/sirupsen/logrus"
	"errors"
	"fmt"
	"path"
)

type DirectoryWatcher struct {
	ticker        *time.Ticker
	watchedFiles  map[string]struct{}
	recordChannel chan *models.RecordCreate
	retryCount    int
}

func NewDirectoryWatcher(scanInterval, retry int) *DirectoryWatcher {
	return &DirectoryWatcher{
		ticker:        time.NewTicker(time.Duration(scanInterval) * time.Second),
		recordChannel: make(chan *models.RecordCreate),
		watchedFiles:  make(map[string]struct{}),
		retryCount:    retry,
	}
}

func (w *DirectoryWatcher) Close() {
	w.ticker.Stop()
	close(w.recordChannel)
}

func (w *DirectoryWatcher) Watch(dir string, parser parser.Parser) <-chan *models.RecordCreate {
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

func (w *DirectoryWatcher) Done(record *models.RecordCreate) {
	if err := w.remove(record); err != nil {
		log.WithField("error", err).WithField("record", record).Infof("error processing record")
	} else {
		log.WithField("record", record).Infof("record sucessfully processed")
	}
}

func (w *DirectoryWatcher) Error(record *models.RecordCreate) {
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

func (w *DirectoryWatcher) add(record *models.RecordCreate) {
	if _, ok := w.watchedFiles[record.PdfPath]; ok {
		return
	}

	w.watchedFiles[record.PdfPath] = struct{}{}

	w.recordChannel <- record
}

func (w *DirectoryWatcher) remove(record *models.RecordCreate) error {
	if _, ok := w.watchedFiles[record.PdfPath]; !ok {
		return errors.New(fmt.Sprintf("record %s not found", record.PdfPath))
	}

	if err := os.Remove(record.PdfPath); err != nil {
		return err
	}

	delete(w.watchedFiles, record.PdfPath)
	return nil
}
