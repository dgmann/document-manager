package main

import (
	"fmt"
	"github.com/dgmann/document-manager/api-client/repository"
	"github.com/dgmann/document-manager/migrator/categories"
	"github.com/dgmann/document-manager/migrator/importer"
	"github.com/dgmann/document-manager/migrator/patients"
	"github.com/dgmann/document-manager/migrator/shared"
	"github.com/gosuri/uiprogress"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"sync"
)

func main() {
	go func() {
		logrus.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	logrus.Info("Importer started")
	config := NewConfig()
	i := importer.NewImporter(config.ApiURL)

	var importData importer.Import
	logrus.WithField("file", config.InputFile).Info("Load records")
	err := importData.Load(config.InputFile)
	if err != nil {
		logrus.WithError(err).Fatal("error opening input file")
		return
	}

	patientsToImport := importData.Patients
	categoriesToImport := importData.Categories
	recordsToImport := importData.Records

	categoryProgressBar := uiprogress.AddBar(len(categoriesToImport)).AppendCompleted().PrependElapsed().PrependFunc(countFunc(len(categoriesToImport)))
	patientProgressBar := uiprogress.AddBar(len(patientsToImport)).AppendCompleted().PrependElapsed().PrependFunc(countFunc(len(patientsToImport)))
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
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := importPatients(i, patientsToImport, patientProgressBar); err != nil {
			patientProgressBar.AppendFunc(logError(err))
			logrus.WithError(err).Error("error importing patients")
		}
	}()

	logrus.WithField("count", len(importData.Records)).Info("Start importing")
	var recordsNotImported []string
	recordProgressBar.AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Errors: %d", len(recordsNotImported))
	})

	imported, notImported := i.ImportRecords(importData.Records)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range imported {
			recordProgressBar.Incr()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range notImported {
			recordProgressBar.Incr()
			logrus.WithError(err).Error("error importing record")
			recordsNotImported = append(recordsNotImported, err.Record.Path)
		}
	}()

	wg.Wait()
	logrus.WithField("errors", len(recordsNotImported)).Info("ImportRecords finished")
	err = shared.WriteLines(recordsNotImported, config.ErrorFile)
	if err != nil {
		logrus.WithError(err).Fatal("error writing output file")
		return
	}
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

func importPatients(i *importer.Importer, patients []*patients.Patient, progressbar *uiprogress.Bar) error {
	logrus.WithField("count", len(patients)).Info("start importing patients")
	for _, patient := range patients {
		p := repository.Patient{
			Id: patient.Id,
		}
		if patient.Name != nil {
			splitted := strings.Split(*patient.Name, ",")
			if len(splitted) == 2 {
				p.LastName = splitted[0]
				p.FirstName = splitted[1]
			}
		}

		if err := i.Import("/patients", p); err != nil {
			return err
		} else {
			progressbar.Incr()
		}
	}
	logrus.Info("patients successfully imported")
	return nil
}

func logError(err error) (func(b *uiprogress.Bar) string) {
	return func(b *uiprogress.Bar) string {
		if err != nil {
			return "Error"
		} else {
			return ""
		}
	}
}

func countFunc(total int) (func(b *uiprogress.Bar) string) {
	return func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%d/%d", b.Current(), total)
	}
}
