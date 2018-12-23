package main

import (
	"fmt"
	"github.com/dgmann/document-manager/migrator/categories"
	"github.com/dgmann/document-manager/migrator/importer"
	"github.com/dgmann/document-manager/migrator/shared"
	"github.com/gosuri/uiprogress"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

func main() {
	go func() {
		logrus.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	logrus.Info("Importer started")
	config := NewConfig()
	i := importer.NewImporter(config.ApiURL, config.RetryCount)

	var importData importer.Import
	logrus.WithField("file", config.InputFile).Info("Load records")
	err := importData.Load(config.InputFile)
	if err != nil {
		logrus.WithError(err).Fatal("error opening input file")
		return
	}

	categoriesToImport := importData.Categories
	recordsToImport := importData.Records
	alreadyImported := make(map[string]importer.ImportableRecord)
	err = shared.LoadFromFile(path.Join(config.DataDirectory, "importedrecords.gob"), alreadyImported)
	if err != nil {
		logrus.WithError(err).Info("no records found which were already imported")
	}

	categoryProgressBar := uiprogress.AddBar(len(categoriesToImport)).AppendCompleted().PrependElapsed().PrependFunc(countFunc(len(categoriesToImport)))
	recordProgressBar := uiprogress.AddBar(len(recordsToImport)).AppendCompleted().PrependElapsed().PrependFunc(countFunc(len(recordsToImport)))
	uiprogress.Start()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := importCategories(i, categoriesToImport, categoryProgressBar); err != nil {
			categoryProgressBar.AppendFunc(logError(err))
			logrus.WithError(err).Error("error importing categories")
		}
	}()

	logrus.WithField("count", len(importData.Records)).Info("Start importing")
	var recordsNotImported []importer.ImportableRecord
	recordProgressBar.AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Errors: %d", len(recordsNotImported))
	})

	records := importer.Difference(recordsToImport, alreadyImported)
	if err := recordProgressBar.Set(len(recordsToImport) - len(records)); err != nil {
		logrus.WithError(err).Warn("error setting progressbar")
	}

	importedRecords := make(map[string]importer.ImportableRecord)
	registerSignals(importedRecords, path.Join(config.DataDirectory, "importedrecords.gob"))
	imported, notImported := i.ImportRecords(records)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range imported {
			recordProgressBar.Incr()
			importedRecords[r.Path] = *r
		}
	}()

	var errorLines []string
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range notImported {
			recordProgressBar.Incr()
			logrus.WithError(err).Error("error importing record ", err.Record.Path)
			recordsNotImported = append(recordsNotImported, *err.Record)
			errorLines = append(errorLines, fmt.Sprintf("%s: %s", err.Record.Path, err.Error()))
		}
	}()

	wg.Wait()
	logrus.WithField("errors", len(recordsNotImported)).Info("ImportRecords finished")

	reimportable := importer.Import{Records: recordsNotImported}
	if err = reimportable.Save(path.Join(config.DataDirectory, "failedrecords.gob")); err != nil {
		logrus.WithError(err).Info("error saving failedrecords.gob")
	}

	err = shared.WriteLines(errorLines, path.Join(config.DataDirectory, "errors.log"))
	if err != nil {
		logrus.WithError(err).Fatal("error writing output file")
		return
	}
}

func registerSignals(importedRecords map[string]importer.ImportableRecord, path string) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-sigs
		defer os.Exit(0)
		logrus.Info("save list of imported records")
		err := shared.SaveToFile(path, importedRecords)
		if err != nil {
			logrus.WithError(err).Info("error saving imported records")
			return
		}
	}()
}

func importCategories(i *importer.Importer, categories []*categories.Category, progressbar *uiprogress.Bar) error {
	logrus.WithField("count", len(categories)).Info("start importing categories")
	for _, category := range categories {
		if err := i.Import("/categories", category); err != nil {
			return err
		} else {
			progressbar.Incr()
		}
	}
	logrus.Info("categories successfully imported")
	return nil
}

func logError(err error) func(b *uiprogress.Bar) string {
	return func(b *uiprogress.Bar) string {
		if err != nil {
			return "Error"
		} else {
			return ""
		}
	}
}

func countFunc(total int) func(b *uiprogress.Bar) string {
	return func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%d/%d", b.Current(), total)
	}
}
