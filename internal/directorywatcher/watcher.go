package directorywatcher

import (
	"context"
	"log/slog"
	"os"

	"github.com/dgmann/document-manager/pkg/client"
	"github.com/dgmann/document-manager/pkg/directory-watcher/parser"
	"github.com/dgmann/document-manager/pkg/directory-watcher/watcher"
	"github.com/dgmann/document-manager/pkg/log"
)

var logger = log.Logger

func Watch(ctx context.Context, config Config) {
	uploader, err := client.NewHTTPClient(config.DestinationUrl, config.Timeout)
	if err != nil {
		logger.Error("error creating API client", log.ErrAttr(err))
		return
	}

	uploadChan := make(chan *watcher.NewRecord)
	defer close(uploadChan)

	// Watch directories
	for _, source := range config.Sources {
		logger.Info("initializing watcher", "directory", source.Directory)
		p := func() parser.Parser {
			if source.Parser == "fax" {
				return &parser.Fax{}
			} else if source.Parser == "generic" {
				return &parser.Generic{
					Sender: source.Sender,
				}
			} else {
				panic("Invalid parser: " + source.Parser)
			}
		}()
		w := watcher.NewDirectoryWatcher(source.Directory, p, config.ScanInterval, config.RetryCount)
		recordsChan := w.Watch()
		go func(recordsChan <-chan *watcher.NewRecord) {
			for {
				select {
				case record := <-recordsChan:
					uploadChan <- record
				case <-ctx.Done():
					w.Close()
				}

			}
		}(recordsChan)
	}
	go func() {
		// Upload files
		logger.Info("starting uploader", "destination", config.DestinationUrl)
		for record := range uploadChan {
			logger := logger.With(slog.String("path", record.PdfPath))
			f, err := os.Open(record.PdfPath)
			if err != nil {
				logger.Error("error opening pdf", log.ErrAttr(err))
				record.Error()
				continue
			}
			record.File = f
			_, err = uploader.Records.Create(record.NewRecord)
			f.Close()
			if err != nil {
				logger.Error("error uploading record", log.ErrAttr(err))
				record.Error()
				continue
			}
			record.Done()
			logger.Info("record uploaded")
		}
	}()
	<-ctx.Done()
}
