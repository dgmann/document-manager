package shared

import "github.com/namsral/flag"

type Config struct {
	Username, Password, Hostname, Instance, RecordDirectory, DbName, ValidationFile, SplittedDirectory, DataDirectory string
}

func NewConfig() Config {
	var username, password, hostname, instance, recordDirectory, dbName, validationFile, splittedDir, dataDirectory string

	flag.StringVar(&username, "db_user", "", "Database Username")
	flag.StringVar(&password, "db_password", "", "Database Password")
	flag.StringVar(&hostname, "db_host", "", "Database Hostname")
	flag.StringVar(&instance, "db_instance", "", "Database Instance name")
	flag.StringVar(&dbName, "db_name", "", "Database name")
	flag.StringVar(&validationFile, "validation_file", "./error.log", "Validation File")

	flag.StringVar(&recordDirectory, "record_dir", "/records", "Record Directory")
	flag.StringVar(&splittedDir, "splitted_dir", "/splitted", "Splitted Records Directory")
	flag.StringVar(&dataDirectory, "data_dir", "/data", "Data Directory")
	flag.Parse()

	return Config{
		Username:          username,
		Password:          password,
		Hostname:          hostname,
		Instance:          instance,
		RecordDirectory:   recordDirectory,
		SplittedDirectory: splittedDir,
		DataDirectory:     dataDirectory,
		DbName:            dbName,
		ValidationFile:    validationFile,
	}
}
