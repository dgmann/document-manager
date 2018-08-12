package main

import (
	"github.com/dgmann/document-manager/migrator/records/databasereader"
	"github.com/dgmann/document-manager/migrator/records/filesystem"
	"github.com/dgmann/document-manager/migrator/shared"
	"github.com/pkg/errors"
	"github.com/dgmann/document-manager/migrator/validator"
	"fmt"
	"strings"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"net/http"
	_ "net/http/pprof"
	"github.com/dgmann/document-manager/migrator/importer"
	"github.com/dgmann/document-manager/api-client/record"
	"strconv"
)

func main() {
	go func() {
		logrus.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	config := NewConfig()
	recordManager := databasereader.NewManager(config.DbName, config.Username, config.Password, config.Hostname, config.Instance)
	err := recordManager.Open()
	if err != nil {
		logrus.WithError(err).Fatal("Error opening connection")
		return
	}
	defer recordManager.Close()

	databaseIndex, filesystemIndex, err := load(config, recordManager)
	if err != nil {
		logrus.WithError(err).Error("error loading data")
	}
	err = filesystemIndex.Save(filepath.Join(config.DataDirectory, "filesystem.gob"))
	if err != nil {
		logrus.WithError(err).Error("error saving filesystemindex to disk")
	}

	resolvable, validationErrors := validator.Validate(filesystemIndex, databaseIndex, recordManager.Manager)
	if validationErrors != nil {
		logrus.WithError(validationErrors).Warn("validation error")
	}

	shared.WriteLines(validationErrors.Messages, config.ValidationFile)
	if askForConfirmation() {
		fmt.Printf("Trying to resolve them\n")
		for _, r := range resolvable {
			r.Resolve()
		}
		return
	}

	var filesToImport importer.ImportableRecordList
	status := record.StatusDone
	for _, r := range filesystemIndex.SubRecords() {
		subrecord := r.SubRecord()
		patId := strconv.Itoa(*subrecord.PatId)
		newRecord := record.CreateRecord{
			ReceivedAt: subrecord.ReceivedAt,
			Date:       *subrecord.Date,
			Status:     &status,
			PatientId:  &patId,
		}
		filesToImport = append(filesToImport, importer.ImportableRecord{
			CreateRecord: newRecord,
			Path:         subrecord.Path,
		})
	}
	err = filesToImport.Save(config.OutputFile)
	if err != nil {
		fmt.Printf("Error creating output %v\n", err)
	} else {
		fmt.Printf("Successfully created file to import\n")
	}
}

func load(config Config, manager *databasereader.Manager) (*databasereader.Index, *filesystem.Index, error) {
	errorChan := make(chan error, 2)
	databaseIndexChan := make(chan *databasereader.Index, 1)
	filesystemIndexChan := make(chan *filesystem.Index, 1)

	go func() {
		index, err := manager.Load()
		if err != nil {
			errorChan <- errors.Wrap(err, "error loading from database")
		}
		databaseIndexChan <- index
	}()

	go func() {
		index, err := filesystem.LoadIndexFromFile(filepath.Join(config.DataDirectory, "filesystem.gob"))
		if err != nil {
			index, err = loadFileSystem(config.RecordDirectory)
		}
		if err != nil {
			errorChan <- errors.Wrap(err, "error loading from filesystem")
		}
		logrus.Info("load sub records")
		index.LoadSubRecords(config.SplittedDirectory)
		filesystemIndexChan <- index
	}()

	databaseIndex := <-databaseIndexChan
	filesystemIndex := <-filesystemIndexChan

	close(errorChan)
	close(databaseIndexChan)
	close(filesystemIndexChan)

	var err error
	for e := range errorChan {
		err = shared.WrapError(err, e.Error())
	}

	return databaseIndex, filesystemIndex, err
}

func loadFileSystem(recordDirectory string) (*filesystem.Index, error) {
	return filesystem.CreateIndex(recordDirectory)
}

func askForConfirmation() bool {
	var s string

	fmt.Printf("Resolve validation errors? (y/N): \n")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}
