package main

import (
	"github.com/dgmann/document-manager/api/client"
	"github.com/dgmann/document-manager/directory-watcher/parser"
	"github.com/dgmann/document-manager/directory-watcher/watcher"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var directory string
var destination string
var pars string
var sender string
var retryCount int
var scanInterval int

func init() {
	flag.StringVar(&directory, "directory", "", "Directory to watch")
	flag.StringVar(&destination, "destination", "", "Upload destination")
	flag.StringVar(&pars, "parser", "generic", "The parser to use to parse the file names")
	flag.StringVar(&sender, "sender", "", "The value to use as a sender")
	flag.IntVar(&retryCount, "retry", 5, "Times to retry uploading a record")
	flag.IntVar(&scanInterval, "scan", 1, "Interval in seconds at which to scan the directory")
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
	uploader := client.NewHttpUploader(destination, time.Second*5)
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
	records := w.Watch(directory, p)
	for record := range records {
		f, err := os.Open(record.PdfPath)
		if err != nil {
			log.WithField("record", record).WithField("error", err).Errorf("error opening pdf")
			w.Error(record)
		}
		record.File = f
		err = uploader.CreateRecord(record.NewRecord)
		f.Close()
		if err != nil {
			log.WithField("record", record).WithField("error", err).Errorf("error uploading record")
			w.Error(record)
		} else {
			w.Done(record)
		}
	}
}
