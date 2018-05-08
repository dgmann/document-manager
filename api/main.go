package main

import (
	"github.com/dgmann/document-manager/api/http"
	"github.com/dgmann/document-manager/api/shared"
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/dgmann/document-manager/api/services"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	recordDir := envOrDefault("RECORD_DIR", "/records")
	dbHost := envOrDefault("DB_HOST", "localhost")
	dbname := envOrDefault("DB_NAME", "manager")
	pdfprocessorUrl := envOrDefault("PDFPROCESSOR_URL", "http://localhost:8181")
	baseUrl := envOrDefault("BASE_URL", "http://localhost:8080")

	session, err := mgo.Dial(dbHost)
	if err != nil {
		log.Errorf("Error connecting to database: %s", err)
	}
	defer session.Close()

	config := &shared.Config{
		Db:              session.DB(dbname),
		RecordDir:       recordDir,
		PdfProcessorUrl: pdfprocessorUrl,
		BaseUrl:         baseUrl,
	}

	services.InitHealthService(dbHost, pdfprocessorUrl)
	factory := http.NewFactory(config)
	http.Run(factory)
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
