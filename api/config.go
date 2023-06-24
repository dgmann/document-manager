package main

import (
	"github.com/dgmann/document-manager/api/internal/datastore"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	RecordDirectory  string
	ArchiveDirectory string
	PdfProcessorUrl  string
	Database         datastore.DatabaseConfig
	MQTTBroker       string
	MQTTClientId     string
	Port             string
}

func ConfigFromEnv() Config {
	recordDir := envOrDefault("RECORD_DIR", "/records")
	archiveDir := envOrDefault("ARCHIVE_DIR", "/archive")
	dbHost := envOrDefault("DB_HOST", "localhost")
	dbPort := envOrDefault("DB_PORT", "27017")
	dbName := envOrDefault("DB_NAME", "manager")
	pdfProcessorUrl := envOrDefault("PDFPROCESSOR_URL", "127.0.0.1:9000")
	mqttBroker := envOrDefault("MQTT_BROKER", "mqtt:1883")
	port := envOrDefault("HTTP_PORT", "80")
	mqttClientId := envOrDefault("MQTT_CLIENT_ID", "backend-api-"+strconv.Itoa(rand.Int()))
	return Config{
		RecordDirectory:  recordDir,
		ArchiveDirectory: archiveDir,
		PdfProcessorUrl:  pdfProcessorUrl,
		Database: datastore.DatabaseConfig{
			Host: strings.TrimSpace(dbHost),
			Port: strings.TrimSpace(dbPort),
			Name: strings.TrimSpace(dbName),
		},
		MQTTBroker:   mqttBroker,
		MQTTClientId: mqttClientId,
		Port:         port,
	}
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
