package main

import (
	"context"
	"fmt"
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

func main() {
	recordDir := envOrDefault("RECORD_DIR", "/records")
	archiveDir := envOrDefault("ARCHIVE_DIR", "/archive")
	dbHost := envOrDefault("DB_HOST", "localhost")
	dbName := envOrDefault("DB_NAME", "manager")
	pdfProcessorUrl := envOrDefault("PDFPROCESSOR_URL", "127.0.0.1:9000")

	if err := ensureTmpDirectory(); err != nil {
		log.Error(fmt.Errorf("error while creating tmp directory: %w", err))
		return
	}

	log.WithFields(log.Fields{"host": dbHost, "database": dbName}).Info("connecting to database")
	client := mongo.NewClient(dbHost, dbName)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.WithError(err).Error("database cannot be reached")
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
	categoryService := mongo.NewCategoryService(client.Categories(), client.Records())

	pdfProcessor, err := grpc.NewPDFProcessor(pdfProcessorUrl, imageService, categoryService)
	if err != nil {
		log.WithError(err).Error("error connecting to pdf processor service")
	}
	eventService := http.NewEventService(imageService)
	tagService := mongo.NewTagService(client.Records())

	srv := http.Server{
		EventService:    eventService,
		ImageService:    imageService,
		TagService:      tagService,
		CategoryService: categoryService,
		ArchiveService:  archiveService,
		RecordService: mongo.NewRecordService(mongo.RecordServiceConfig{
			Records: client.Records(),
			Events:  eventService,
			Images:  imageService,
			Pdfs:    archiveService,
		}),
		PdfProcessor: pdfProcessor,
		Healthchecker: map[string]app.Checkable{
			"database":       client,
			"pdfProcessor":   pdfProcessor,
			"archiveStorage": archiveService,
			"recordStorage":  imageService,
		},
		StatisticProviders: map[string]app.StatisticProvider{
			"archiveStorage": archiveService,
		},
	}

	log.Info("server startup completed")
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

func ensureTmpDirectory() error {
	if _, err := os.Stat(os.TempDir()); os.IsNotExist(err) {
		if err := os.Mkdir(os.TempDir(), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
