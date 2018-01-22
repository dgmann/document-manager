package main

import (
	"github.com/dgmann/document-manager-api/http"
	"github.com/dgmann/document-manager-api/repositories"
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	recordDir := envOrDefault("RECORD_DIR", "./records")
	host := envOrDefault("DB_HOST", "localhost")
	dbname := envOrDefault("DB_NAME", "manager")

	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(dbname).C("records")
	records := repositories.NewRecordRepository(c)
	http.Run(records, recordDir)
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
