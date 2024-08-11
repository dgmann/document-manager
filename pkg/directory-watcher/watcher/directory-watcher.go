package watcher

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/dgmann/document-manager/pkg/api"
	"github.com/dgmann/document-manager/pkg/directory-watcher/parser"
	"github.com/dgmann/document-manager/pkg/log"
)

type NewRecord struct {
	*api.NewRecord
	PdfPath      string
	RetryCounter int
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
			files, err := os.ReadDir(dir)
			if err != nil {
				logger.Error("error reading directory", log.ErrAttr(err), slog.String("directory", dir))
				continue
			}

			for _, f := range files {
				if f.IsDir() || path.Ext(f.Name()) != ".pdf" {
					continue
				}

				parsed := parser.Parse(f.Name())
				record := &NewRecord{NewRecord: parsed, PdfPath: f.Name(), RetryCounter: 0}
				record.PdfPath = path.Join(dir, f.Name())
				w.add(record)
			}
		}
	}()
	return w.recordChannel
}

func (w *DirectoryWatcher) Done(record *NewRecord) {
	if err := w.remove(record); err != nil {
		logger.With(slog.String("file", record.PdfPath)).Info("could not remove processed record", log.ErrAttr(err))
	} else {
		logger.With(slog.String("file", record.PdfPath)).Info("record sucessfully processed")
	}
}

func (w *DirectoryWatcher) Error(record *NewRecord) {
	record.RetryCounter++
	if record.RetryCounter <= w.retryCount {
		go func(record *NewRecord) {
			time.Sleep(2 * time.Second)
			w.recordChannel <- record
			logger.With(slog.String("file", record.PdfPath)).Info("requeue record")
		}(record)
	} else {
		logger.With(slog.String("file", record.PdfPath), slog.Int("retryCount", w.retryCount)).Info("retry counter exceeded")
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
		return fmt.Errorf("record %s not found", record.PdfPath)
	}

	if err := os.Remove(record.PdfPath); err != nil {
		return err
	}

	delete(w.watchedFiles, record.PdfPath)
	return nil
}
