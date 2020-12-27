package main

import (
	"os"
	"strings"

	"github.com/dgmann/document-manager/api/datastore"
)

type Config struct {
	RecordDirectory  string
	ArchiveDirectory string
	PdfProcessorUrl  string
	Database         datastore.DatabaseConfig
}

func ConfigFromEnv() Config {
	recordDir := envOrDefault("RECORD_DIR", "/records")
	archiveDir := envOrDefault("ARCHIVE_DIR", "/archive")
	dbHost := envOrDefault("DB_HOST", "localhost")
	dbPort := envOrDefault("DB_PORT", "27017")
	dbName := envOrDefault("DB_NAME", "manager")
	pdfProcessorUrl := envOrDefault("PDFPROCESSOR_URL", "127.0.0.1:9000")
	return Config{
		RecordDirectory:  recordDir,
		ArchiveDirectory: archiveDir,
		PdfProcessorUrl:  pdfProcessorUrl,
		Database: datastore.DatabaseConfig{
			Host: strings.TrimSpace(dbHost),
			Port: strings.TrimSpace(dbPort),
			Name: strings.TrimSpace(dbName),
		},
	}
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
