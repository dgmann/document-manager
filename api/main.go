package main

import (
	"github.com/dgmann/document-manager/api/http"
	"github.com/dgmann/document-manager/api/pdf"
	"github.com/dgmann/document-manager/api/repositories"
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
	recordDir := envOrDefault("RECORD_DIR", "C:\\Users\\David\\Desktop\\Images")
	dbHost := envOrDefault("DB_HOST", "localhost")
	dbname := envOrDefault("DB_NAME", "manager")
	pdfprocessorUrl := envOrDefault("PDFPROCESSOR_URL", "http://localhost:8181")

	session, err := mgo.Dial(dbHost)
	if err != nil {
		log.Errorf("Error connecting to database: %s", err)
	}
	defer session.Close()
	c := session.DB(dbname).C("records")
	images := repositories.NewFileSystemImageRepository(recordDir)
	app := shared.App{
		Records:      repositories.NewRecordRepository(c, images),
		Images:       images,
		Tags:         repositories.NewTagRepository(c),
		Categories:   repositories.NewCategoryRepository(),
		PDFProcessor: pdf.NewPDFProcessor(pdfprocessorUrl),
	}
	services.InitHealthService(dbHost, pdfprocessorUrl)
	http.Run(&app)
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}