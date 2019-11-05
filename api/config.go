package main

import (
	"os"
)

type Config struct {
	RecordDirectory  string
	ArchiveDirectory string
	PdfProcessorUrl  string
	Database         struct {
		Host string
		Name string
	}
}

func ConfigFromEnv() Config {
	recordDir := envOrDefault("RECORD_DIR", "/records")
	archiveDir := envOrDefault("ARCHIVE_DIR", "/archive")
	dbHost := envOrDefault("DB_HOST", "localhost")
	dbName := envOrDefault("DB_NAME", "manager")
	pdfProcessorUrl := envOrDefault("PDFPROCESSOR_URL", "127.0.0.1:9000")
	return Config{
		RecordDirectory:  recordDir,
		ArchiveDirectory: archiveDir,
		PdfProcessorUrl:  pdfProcessorUrl,
		Database: struct {
			Host string
			Name string
		}{Host: dbHost, Name: dbName},
	}
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
