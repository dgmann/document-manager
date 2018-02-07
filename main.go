package main

import (
	"github.com/dgmann/document-manager-api/http"
	"github.com/dgmann/document-manager-api/repositories"
	"github.com/dgmann/document-manager-api/shared"
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/dgmann/document-manager-api/services"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	recordDir := envOrDefault("RECORD_DIR", "C:\\Users\\David\\Desktop\\Images")
	host := envOrDefault("DB_HOST", "localhost")
	dbname := envOrDefault("DB_NAME", "manager")

	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(dbname).C("records")
	images := repositories.NewFileSystemImageRepository(recordDir)
	app := shared.App{
		Records: repositories.NewRecordRepository(c, images),
		Images: images,
		PDFProcessor: services.NewPDFProcessor("http://10.0.0.38:8181"),
	}
	http.Run(&app)
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
