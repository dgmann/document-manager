package main

import (
	"context"
	"github.com/Shopify/logrus-bugsnag"
	"github.com/bugsnag/bugsnag-go"
	"github.com/dgmann/document-manager/api/http"
	"github.com/dgmann/document-manager/api/repositories/record"
	"github.com/dgmann/document-manager/api/services"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var bugsnagConfig bugsnag.Configuration

func init() {
	apiKey := envOrDefault("BUGSNAG_API_KEY", "")
	stage := envOrDefault("RELEASE_STAGE", "production")
	bugsnagConfig = bugsnag.Configuration{
		// Your Bugsnag project API key
		APIKey:       apiKey,
		ReleaseStage: stage,
		// The import paths for the Go packages
		// containing your source files
		ProjectPackages: []string{"main", "github.com/dgmann/document-manager/api"},
		PanicHandler:    func() {},
	}
	bugsnag.Configure(bugsnagConfig)

	hook, err := logrus_bugsnag.NewBugsnagHook()
	if err != nil {
		log.Error("Error registering bugsnag hook")
	}
	log.StandardLogger().Hooks.Add(hook)
}

func main() {
	recordDir := envOrDefault("RECORD_DIR", "/records")
	archiveDir := envOrDefault("ARCHIVE_DIR", "/archive")
	dbHost := envOrDefault("DB_HOST", "localhost")
	dbname := envOrDefault("DB_NAME", "manager")
	pdfprocessorUrl := envOrDefault("PDFPROCESSOR_URL", "127.0.0.1:9000")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, "mongodb://"+dbHost)
	if err != nil {
		log.WithError(err).Error("error creating database client")
		os.Exit(1)
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.WithError(err).Error("error connecting to database")
		os.Exit(1)
		return
	}

	if err := record.CreateIndexes(context.Background(), client.Database(dbname).Collection("records")); err != nil {
		log.WithError(err).Error("error setting indices")
	}
	config := &Config{
		Db:              client.Database(dbname),
		RecordDir:       recordDir,
		PDFDir:          archiveDir,
		PdfProcessorUrl: pdfprocessorUrl,
		Bugsnag:         bugsnagConfig,
	}

	services.InitHealthService(dbHost, pdfprocessorUrl)
	factory := NewFactory(config)
	http.Run(factory, bugsnagConfig)
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
