package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/dgmann/document-manager/api/pkg/client"
	"github.com/dgmann/document-manager/directory-watcher/parser"
	"github.com/dgmann/document-manager/directory-watcher/watcher"
	"github.com/dgmann/document-manager/pkg/log"
	"github.com/namsral/flag"
)

var logger = log.Logger

var directory string
var destination string
var pars string
var sender string
var retryCount int
var scanInterval int
var timeout int

func init() {
	flag.StringVar(&directory, "directory", "", "Directory to watch")
	flag.StringVar(&destination, "destination", "", "Upload destination")
	flag.StringVar(&pars, "parser", "generic", "The parser to use to parse the file names")
	flag.StringVar(&sender, "sender", "", "The value to use as a sender")
	flag.IntVar(&retryCount, "retry", 5, "Times to retry uploading a record")
	flag.IntVar(&scanInterval, "scan", 1, "Interval in seconds at which to scan the directory")
	flag.IntVar(&timeout, "timeout", 60, "timeout in seconds. Default: 60")
	flag.Parse()
	if len(directory) == 0 {
		panic("Invalid directory")
	}

	if len(destination) == 0 {
		panic("Invalid destination")
	}
}

func main() {
	w := watcher.NewDirectoryWatcher(scanInterval, retryCount)
	uploader, err := client.NewHTTPClient(destination, time.Second*time.Duration(timeout))
	if err != nil {
		logger.Error("error creating API client", log.ErrAttr(err))
		os.Exit(1)
		return
	}
	var p parser.Parser
	if pars == "fax" {
		p = &parser.Fax{}
	} else if pars == "generic" {
		p = &parser.Generic{
			Sender: sender,
		}
	} else {
		panic("Invalid parser: " + pars)
	}
	logger.Info("start watching directory", "directory", directory)
	records := w.Watch(directory, p)
	for record := range records {
		logger := logger.With(slog.String("path", record.PdfPath))
		f, err := os.Open(record.PdfPath)
		if err != nil {
			logger.Error("error opening pdf", log.ErrAttr(err))
			w.Error(record)
			continue
		}
		record.File = f
		_, err = uploader.Records.Create(record.NewRecord)
		f.Close()
		if err != nil {
			logger.Error("error uploading record", log.ErrAttr(err))
			w.Error(record)
		} else {
			w.Done(record)
		}
	}
}
