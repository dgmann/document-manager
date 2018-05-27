package main

import (
	"net/url"
	"github.com/dgmann/document-manager/migrator/databasereader"
	"flag"
)

var username, password, hostname, instance string

func init() {
	flag.StringVar(&username, "u", "", "Database username")
	flag.StringVar(&password, "p", "", "Database password")
	flag.StringVar(&hostname, "h", "", "Database hostname")
	flag.StringVar(&instance, "i", "SQLExpress", "Database instance name")
	flag.Parse()
}

func main() {
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

	records, err := manager.Load()
	if err != nil {
		println("Error loading from database", err.Error())
	}
	println("Patient count: ", len(records))
}
