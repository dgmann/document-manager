package main

import (
	"github.com/dgmann/document-manager/migrator/importer"
	"github.com/sirupsen/logrus"
	"net/http"
	"github.com/dgmann/document-manager/migrator/shared"
)

func main() {
	go func() {
		logrus.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	logrus.Info("Importer started")
	config := NewConfig()
	i := importer.NewImporter(config.ApiURL)

	var importableRecords importer.ImportableRecordList
	logrus.WithField("file", config.InputFile).Info("Load records")
	err := importableRecords.Load(config.InputFile)
	if err != nil {
		logrus.WithError(err).Fatal("error opening input file")
		return
	}
	logrus.WithField("count", len(importableRecords)).Info("Start importing")
	notImported := i.Import(importableRecords)
	logrus.WithField("unsuccessful", len(notImported)).Info("Import finished")
	err = shared.WriteLines(notImported, config.ErrorFile)
	if err != nil {
		logrus.WithError(err).Fatal("error writing output file")
		return
	}
}
