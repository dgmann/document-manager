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

	config := NewConfig()
	i := importer.NewImporter(config.ApiURL)

	var importableRecords importer.ImportableRecordList
	err := importableRecords.Load(config.InputFile)
	if err != nil {
		logrus.WithError(err).Fatal("error opening input file")
		return
	}
	notImported := i.Import(importableRecords)
	err = shared.WriteLines(notImported, config.OutputFile)
	if err != nil {
		logrus.WithError(err).Fatal("error writing output file")
		return
	}
}
