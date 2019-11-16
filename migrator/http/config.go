package http

import (
	"github.com/namsral/flag"
	"path/filepath"
)

type Config struct {
	Username, Password, Hostname, Instance, RecordDirectory, DbName, DataDirectory, OutputFile, ApiURL string
	RetryCount                                                                                         int
}

func NewConfig() Config {
	var username, password, hostname, instance, recordDirectory, dbName, validationFile, dataDirectory, apiURL string
	var retryCount int

	flag.StringVar(&username, "db_user", "", "Database Username")
	flag.StringVar(&password, "db_password", "", "Database Password")
	flag.StringVar(&hostname, "db_host", "", "Database Hostname")
	flag.StringVar(&instance, "db_instance", "", "Database Instance name")
	flag.StringVar(&dbName, "db_name", "", "Database name")
	flag.StringVar(&validationFile, "validation_file", "/data/error.log", "Validation File")

	flag.StringVar(&recordDirectory, "record_dir", "/records", "Record Directory")
	flag.StringVar(&dataDirectory, "data_dir", "/data", "Data Directory")

	flag.StringVar(&apiURL, "api_url", "http://api/api", "The URL of the API")
	flag.IntVar(&retryCount, "retry_counter", 3, "Number of times to retry uploading a file after a failure")

	flag.Parse()

	return Config{
		Username:        username,
		Password:        password,
		Hostname:        hostname,
		Instance:        instance,
		RecordDirectory: recordDirectory,
		DataDirectory:   dataDirectory,
		DbName:          dbName,
		OutputFile:      filepath.Join(dataDirectory, "output.gob"),
		ApiURL:          apiURL,
		RetryCount:      retryCount,
	}
}
