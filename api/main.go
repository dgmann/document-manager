package main

import (
	"github.com/bugsnag/bugsnag-go"
	"github.com/dgmann/document-manager/api/http"
	"github.com/dgmann/document-manager/api/services"
	"github.com/dgmann/document-manager/api/shared"
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var bugsnagConfig bugsnag.Configuration

func main() {
	recordDir := envOrDefault("RECORD_DIR", "/records")
	archiveDir := envOrDefault("ARCHIVE_DIR", "/archive")
	dbHost := envOrDefault("DB_HOST", "localhost")
	dbname := envOrDefault("DB_NAME", "manager")
	pdfprocessorUrl := envOrDefault("PDFPROCESSOR_URL", "127.0.0.1:9000")

	session, err := mgo.Dial(dbHost)
	if err != nil {
		log.Errorf("Error connecting to database: %s", err)
	}
	defer session.Close()

	config := &shared.Config{
		Db:              session.DB(dbname),
		RecordDir:       recordDir,
		PDFDir:          archiveDir,
		PdfProcessorUrl: pdfprocessorUrl,
		Bugsnag:         bugsnagConfig,
	}

	services.InitHealthService(dbHost, pdfprocessorUrl)
	factory := http.NewFactory(config)
	http.Run(factory, config)
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
