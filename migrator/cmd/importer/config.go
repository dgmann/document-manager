package main

import "github.com/namsral/flag"

type Config struct {
	ApiURL        string
	InputFile     string
	DataDirectory string
}

func NewConfig() Config {
	var apiURL, inputFile, dataDir string

	flag.StringVar(&inputFile, "input", "", "File containing file paths to import")
	flag.StringVar(&dataDir, "data_dir", "/data", "Directory to store data")
	flag.StringVar(&apiURL, "api_url", "http://localhost", "The URL of the API")
	flag.Parse()

	return Config{
		ApiURL:        apiURL,
		InputFile:     inputFile,
		DataDirectory: dataDir,
	}
}
