package main

import (
	"context"
	"github.com/Shopify/logrus-bugsnag"
	"github.com/bugsnag/bugsnag-go"
	"github.com/dgmann/document-manager/api/app/filesystem"
	"github.com/dgmann/document-manager/api/app/grpc"
	"github.com/dgmann/document-manager/api/app/http"
	"github.com/dgmann/document-manager/api/app/mongo"
	"github.com/dgmann/document-manager/api/services"
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
	client := mongo.NewClient(dbHost, dbname)
	if err := client.Connect(ctx); err != nil {
		log.WithError(err).Error("error connecting to mongodb")
	}

	if err := client.CreateIndexes(context.Background()); err != nil {
		log.WithError(err).Error("error setting indices")
	}

	imageService, err := filesystem.NewImageService(recordDir)
	if err != nil {
		log.WithError(err).Error("error creating image service")
		return
	}
	archiveService, err := filesystem.NewArchiveService(archiveDir)
	if err != nil {
		log.WithError(err).Error("error creating archive service")
		return
	}
	pdfProcessor, err := grpc.NewPDFProcessor(pdfprocessorUrl)
	if err != nil {
		log.WithError(err).Error("error connecting to pdf processor service")
		return
	}
	eventService := http.NewEventService(imageService)
	tagService := mongo.NewTagService(client.Records())

	srv := http.Server{
		EventService:    eventService,
		ImageService:    imageService,
		TagService:      tagService,
		CategoryService: mongo.NewCategoryService(client.Records(), client.Categories()),
		ArchiveService:  archiveService,
		Bug:             bugsnagConfig,
		RecordService: mongo.NewRecordService(mongo.RecordServiceConfig{
			Records: client.Records(),
			Events:  eventService,
		}),
		PdfProcessor: pdfProcessor,
	}

	services.InitHealthService(dbHost, pdfprocessorUrl)
	srv.Run()
}

func envOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
