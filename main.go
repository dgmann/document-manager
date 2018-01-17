package main

import (
	"github.com/dgmann/document-manager-api/http"
	"github.com/dgmann/document-manager-api/repositories"
	"github.com/dgmann/document-manager-api/services"
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
	dsn, ok := os.LookupEnv("DB_CONN")
	if !ok {
		log.Fatal("Database connection not set")
	}
	db := services.NewPostgresConnection(dsn)
	defer db.Close()

	services.MigratePostgres(db.DB)

	records := repositories.NewRecordRepository(services.NewRecordService(db))
	http.Run(records, recordDir)
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
