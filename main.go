package main

import (
	"fmt"
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
	host := envOrDefault("DB_HOST", "localhost")
	port := envOrDefault("DB_PORT", "5432")
	user := envOrDefault("DB_USER", "postgres")
	password := envOrDefault("DB_PASSWORD", "postgres")
	dbname := envOrDefault("DB_NAME", "manager")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
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
