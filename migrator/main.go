package main

import (
	"net/url"
	"github.com/dgmann/document-manager/migrator/databasereader"
	"flag"
	"github.com/dgmann/document-manager/migrator/filesystem"
	"github.com/dgmann/document-manager/migrator/shared"
	"github.com/pkg/errors"
	"github.com/dgmann/document-manager/migrator/validator"
	"fmt"
)

func main() {
	databaseIndex, filesystemIndex, err := load()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	err = validator.Validate(filesystemIndex, databaseIndex)
	if err != nil {
		fmt.Printf("Validation error: %s", err)
	}
}

func load() (*databasereader.Index, *filesystem.Index, error) {
	errorChan := make(chan error, 2)
	databaseIndexChan := make(chan *databasereader.Index, 1)
	filesystemIndexChan := make(chan *filesystem.Index, 1)

	go func() {
		index, err := loadDatabaseRecords()
		if err != nil {
			errorChan <- errors.Wrap(err, "error loading from database")
		}
		databaseIndexChan <- index
	}()

	go func() {
		index, err := loadFileSystem()
		if err != nil {
			errorChan <- errors.Wrap(err, "error loading from filesystem")
		}
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

func loadDatabaseRecords() (*databasereader.Index, error) {
	var username, password, hostname, instance string

	flag.StringVar(&username, "u", "", "Database username")
	flag.StringVar(&password, "p", "", "Database password")
	flag.StringVar(&hostname, "h", "", "Database hostname")
	flag.StringVar(&instance, "i", "SQLExpress", "Database instance name")
	flag.Parse()

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(username, password),
		Host:   hostname,
		Path:   instance,
	}

	manager := databasereader.Manager{}
	err := manager.Open(u.String())
	if err != nil {
		println("Error opening connection: ", err)
	}
	defer manager.Close()

	return manager.Load()
}

func loadFileSystem() (*filesystem.Index, error) {
	var recordDirectory string

	flag.StringVar(&recordDirectory, "d", "", "Record Directory")
	flag.Parse()

	return filesystem.CreateIndex(recordDirectory)
}
