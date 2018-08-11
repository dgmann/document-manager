package main

import (
	"github.com/dgmann/document-manager/migrator/importer"
	"os"
	"github.com/sirupsen/logrus"
	"bufio"
	"net/http"
	"github.com/dgmann/document-manager/migrator/shared"
)

func main() {
	go func() {
		logrus.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	config := NewConfig()
	i := importer.NewImporter(config.ApiURL)

	paths, err := readPathsFromFile(config.InputFile)
	if err != nil {
		logrus.WithError(err).Fatal("error opening input file")
		return
	}
	notImported := i.Import(paths)
	err = shared.WriteLines(notImported, config.OutputFile)
	if err != nil {
		logrus.WithError(err).Fatal("error writing output file")
		return
	}
}

func readPathsFromFile(path string) ([]string, error) {
	input, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	var lines []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
