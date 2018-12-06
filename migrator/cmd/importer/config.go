package main

import "github.com/namsral/flag"

type Config struct {
	ApiURL        string
	InputFile     string
	DataDirectory string
	RetryCount    int
}

func NewConfig() Config {
	var apiURL, inputFile, dataDir string
	var retryCount int

	flag.StringVar(&inputFile, "input", "", "File containing file paths to import")
	flag.StringVar(&dataDir, "data_dir", "/data", "Directory to store data")
	flag.StringVar(&apiURL, "api_url", "http://localhost", "The URL of the API")
	flag.IntVar(&retryCount, "retry_counter", 3, "Number of times to retry uploading a file after a failure")
	flag.Parse()

	return Config{
		ApiURL:        apiURL,
		InputFile:     inputFile,
		DataDirectory: dataDir,
		RetryCount:    retryCount,
	}
}
