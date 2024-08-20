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
	watcher      *DirectoryWatcher
}

func (r *NewRecord) Done() {
	r.watcher.Done(r)
}

func (r *NewRecord) Error() {
	r.watcher.Error(r)
}

type DirectoryWatcher struct {
	ticker        *time.Ticker
	watchedFiles  map[string]struct{}
	recordChannel chan *NewRecord
	retryCount    int
	parser        parser.Parser
	directory     string
}

func NewDirectoryWatcher(dir string, parser parser.Parser, scanInterval time.Duration, retry int) *DirectoryWatcher {
	return &DirectoryWatcher{
		ticker:        time.NewTicker(scanInterval),
		recordChannel: make(chan *NewRecord),
		watchedFiles:  make(map[string]struct{}),
		retryCount:    retry,
		parser:        parser,
		directory:     dir,
	}
}

func (w *DirectoryWatcher) Close() {
	w.ticker.Stop()
	close(w.recordChannel)
}

func (w *DirectoryWatcher) Watch() <-chan *NewRecord {
	go func() {
		for range w.ticker.C {
			files, err := os.ReadDir(w.directory)
			if err != nil {
				logger.Error("error reading directory", log.ErrAttr(err), slog.String("directory", w.directory))
				continue
			}

			for _, f := range files {
				if f.IsDir() || path.Ext(f.Name()) != ".pdf" {
					continue
				}

				parsed := w.parser.Parse(f.Name())
				record := &NewRecord{NewRecord: parsed, PdfPath: f.Name(), RetryCounter: 0, watcher: w}
				record.PdfPath = path.Join(w.directory, f.Name())
				w.add(record)
			}
		}
	}()
	return w.recordChannel
}

func (w *DirectoryWatcher) Done(record *NewRecord) {
	if err := w.remove(record); err != nil {
		logger.With(slog.String("file", record.PdfPath)).Info("could not remove processed record", log.ErrAttr(err))
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
