package main

import (
	"github.com/dgmann/document-manager/migrator/importer"
	"github.com/sirupsen/logrus"
	"net/http"
	"github.com/dgmann/document-manager/migrator/shared"
	"sync"
	"github.com/dgmann/document-manager/api-client/repository"
	"strings"
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		logrus.WithField("count", len(importData.Categories)).Info("start importing categories")
		for _, category := range importData.Categories {
			i.Import("/categories", category)
		}
		logrus.Info("categories successfully imported")
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		logrus.WithField("count", len(importData.Patients)).Info("start importing patients")
		for _, patient := range importData.Patients {
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
			i.Import("/patients", p)
		}
		logrus.Info("patients successfully imported")
		wg.Done()
	}()

	logrus.WithField("count", len(importData.Records)).Info("Start importing")
	notImported := i.ImportRecords(importData.Records)

	wg.Wait()
	logrus.WithField("unsuccessful", len(notImported)).Info("ImportRecords finished")
	err = shared.WriteLines(notImported, config.ErrorFile)
	if err != nil {
		logrus.WithError(err).Fatal("error writing output file")
		return
	}
}
