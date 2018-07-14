package main

import (
	"net/url"
	"github.com/dgmann/document-manager/migrator/databasereader"
	"github.com/dgmann/document-manager/migrator/filesystem"
	"github.com/dgmann/document-manager/migrator/shared"
	"github.com/pkg/errors"
	"github.com/dgmann/document-manager/migrator/validator"
	"fmt"
	"os"
	"bufio"
)

func main() {
	config := shared.NewConfig()
	databaseIndex, filesystemIndex, err := load(config)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	validationErrors := validator.Validate(filesystemIndex, databaseIndex)
	if validationErrors != nil {
		fmt.Printf("Validation error: %s", validationErrors.Error())
	}
	writeLines(validationErrors.Messages, config.ValidationFile)
}

func load(config shared.Config) (*databasereader.Index, *filesystem.Index, error) {
	errorChan := make(chan error, 2)
	databaseIndexChan := make(chan *databasereader.Index, 1)
	filesystemIndexChan := make(chan *filesystem.Index, 1)

	go func() {
		index, err := loadDatabaseRecords(config.Username, config.Password, config.Hostname, config.Instance, config.DbName)
		if err != nil {
			errorChan <- errors.Wrap(err, "error loading from database")
		}
		databaseIndexChan <- index
	}()

	go func() {
		index, err := loadFileSystem(config.RecordDirectory)
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

func loadDatabaseRecords(username, password, hostname, instance, databasename string) (*databasereader.Index, error) {
	query := url.Values{}
	query.Add("database", databasename)
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(username, password),
		Host:     hostname,
		Path:     instance,
		RawQuery: query.Encode(),
	}

	manager := databasereader.Manager{}
	err := manager.Open(u.String())
	if err != nil {
		println("Error opening connection: ", err)
	}
	defer manager.Close()

	return manager.Load()
}

func loadFileSystem(recordDirectory string) (*filesystem.Index, error) {
	return filesystem.CreateIndex(recordDirectory)
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
