package main

import (
	"context"
	"github.com/Shopify/logrus-bugsnag"
	"github.com/bugsnag/bugsnag-go"
	"github.com/dgmann/document-manager/api/app"
	"github.com/dgmann/document-manager/api/app/filesystem"
	"github.com/dgmann/document-manager/api/app/grpc"
	"github.com/dgmann/document-manager/api/app/http"
	"github.com/dgmann/document-manager/api/app/mongo"
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
	dbName := envOrDefault("DB_NAME", "manager")
	pdfProcessorUrl := envOrDefault("PDFPROCESSOR_URL", "127.0.0.1:9000")

	client := mongo.NewClient(dbHost, dbName)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.WithError(err).Error("error connecting to mongodb")
	}

	if err := client.CreateIndexes(context.Background()); err != nil {
		log.WithError(err).Error("error setting indices")
	}

	imageService, err := filesystem.NewImageService(recordDir)
	if err != nil {
		log.WithError(err).Error("error creating image service")
	}
	archiveService, err := filesystem.NewArchiveService(archiveDir)
	if err != nil {
		log.WithError(err).Error("error creating archive service")
	}
	pdfProcessor, err := grpc.NewPDFProcessor(pdfProcessorUrl)
	if err != nil {
		log.WithError(err).Error("error connecting to pdf processor service")
	}
	eventService := http.NewEventService(imageService)
	tagService := mongo.NewTagService(client.Records())

	srv := http.Server{
		EventService:    eventService,
		ImageService:    imageService,
		TagService:      tagService,
		CategoryService: mongo.NewCategoryService(client.Categories(), client.Records()),
		ArchiveService:  archiveService,
		Bug:             bugsnagConfig,
		RecordService: mongo.NewRecordService(mongo.RecordServiceConfig{
			Records: client.Records(),
			Events:  eventService,
		}),
		PdfProcessor: pdfProcessor,
		Healthchecker: map[string]app.HealthChecker{
			"database":       client,
			"pdfProcessor":   pdfProcessor,
			"archiveStorage": archiveService,
			"recordStorage":  imageService,
		},
	}

	if err := srv.Run(); err != nil {
		log.WithError(err).Error("error starting http server")
	}
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
