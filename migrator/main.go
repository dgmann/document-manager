package main

import (
	"github.com/dgmann/document-manager/migrator/records/databasereader"
	"github.com/dgmann/document-manager/migrator/records/filesystem"
	"github.com/dgmann/document-manager/migrator/shared"
	"github.com/pkg/errors"
	"github.com/dgmann/document-manager/migrator/validator"
	"fmt"
	"os"
	"bufio"
	"strings"
)

func main() {
	config := shared.NewConfig()
	recordManager := databasereader.NewManager(config)
	err := recordManager.Open()
	if err != nil {
		println("Error opening connection: ", err.Error())
		return
	}
	defer recordManager.Close()

	databaseIndex, filesystemIndex, err := load(config, recordManager)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	resolvable, validationErrors := validator.Validate(filesystemIndex, databaseIndex, recordManager.Manager)
	if validationErrors != nil {
		fmt.Printf("Validation error: %s\n", validationErrors.Error())
	}

	writeLines(validationErrors.Messages, config.ValidationFile)
	if !askForConfirmation() {
		fmt.Printf("Aborted")
		return
	}
	println(len(resolvable))
}

func load(config shared.Config, manager *databasereader.Manager) (*databasereader.Index, *filesystem.Index, error) {
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

func askForConfirmation() bool {
	var s string

	fmt.Printf("Resolve validation errors? (y/N): ")
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