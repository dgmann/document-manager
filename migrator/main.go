package main

import (
	"net/url"
	"github.com/dgmann/document-manager/migrator/databasereader"
	"os"
	"flag"
	"github.com/dgmann/document-manager/migrator/filesystem"
	"fmt"
)

func main() {
	argsWithoutProg := os.Args[1:]
	command := argsWithoutProg[len(argsWithoutProg)-1]
	if command == "database" {
		loadDatabaseRecords()
	} else if command == "filesystem" {
		loadFileSystem()
	} else {
		println("Please specify a command. Available commands: database, filesystem")
	}
}

func loadDatabaseRecords() {
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

	index, err := manager.Load()
	if err != nil {
		fmt.Println("Error loading from database: ", err)
		return
	}
	println("Patient count: ", index.GetTotalPatientCount())
	println("Record count: ", index.GetTotalPatientCount())
}

func loadFileSystem() {
	var recordDirectory string

	flag.StringVar(&recordDirectory, "d", "", "Record Directory")
	flag.Parse()

	index, err := filesystem.CreateIndex(recordDirectory)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer index.Destroy()

	fmt.Println(index)
	println("Patient count: ", index.GetTotalPatientCount())
	println("Record count: ", index.GetTotalCategorizableCount())
}
