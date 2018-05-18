package main

import (
	"os"
	"github.com/dgmann/document-manager/m1-adapter/m1"
)

var dsn string

func init() {
	dsn = getEnv("DSN", "")
	if len(dsn) == 0 {
		panic("invalid connection string: " + dsn)
	}
}

func main() {
	adapter := m1.NewDatabaseAdapter(dsn)
	err := adapter.Connect()
	if err != nil {
		println(err)
	}
	defer adapter.Close()
	pat, err := adapter.GetPatient("3")
	if err != nil {
		println(err)
	}
	println(pat)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
