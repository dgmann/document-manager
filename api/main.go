package main

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/api/datastore/mongo"
	"github.com/dgmann/document-manager/api/event"
	"github.com/dgmann/document-manager/api/http"
	"github.com/dgmann/document-manager/api/pdf/grpc"
	"github.com/dgmann/document-manager/api/status"
	"github.com/dgmann/document-manager/api/storage/filesystem"
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
	config := ConfigFromEnv()
	if err := ensureTmpDirectory(); err != nil {
		log.Error(fmt.Errorf("error while creating tmp directory: %w", err))
		return
	}

	log.WithFields(log.Fields{"host": config.Database.Host, "database": config.Database.Name}).Info("connecting to database")
	client := mongo.NewClient(config.Database.Host, config.Database.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.WithError(err).Error("database cannot be reached")
	}

	if err := client.CreateIndexes(context.Background()); err != nil {
		log.WithError(err).Error("error setting indices")
	}

	imageService, err := filesystem.NewImageService(config.RecordDirectory)
	if err != nil {
		log.WithError(err).Error("error creating image service")
	}
	archiveService, err := filesystem.NewArchiveService(config.ArchiveDirectory)
	if err != nil {
		log.WithError(err).Error("error creating archive service")
	}
	categoryService := mongo.NewCategoryService(mongo.NewCollection(client.Categories()), mongo.NewCollection(client.Records()))

	pdfProcessor, err := grpc.NewPDFProcessor(config.PdfProcessorUrl, imageService, categoryService)
	if err != nil {
		log.WithError(err).Error("error connecting to pdf processor service")
	}
	eventService := event.NewWebsocketEventService(imageService)
	tagService := mongo.NewTagService(client.Records())

	srv := http.Server{
		EventService:    eventService,
		ImageService:    imageService,
		TagService:      tagService,
		CategoryService: categoryService,
		ArchiveService:  archiveService,
		RecordService: mongo.NewRecordService(mongo.RecordServiceConfig{
			Records: mongo.NewCollection(client.Records()),
			Events:  eventService,
			Images:  imageService,
			Pdfs:    archiveService,
		}),
		PdfProcessor: pdfProcessor,
		HealthService: status.NewHealthService(status.HealthCheckables{
			"database":       client,
			"pdfProcessor":   pdfProcessor,
			"archiveStorage": archiveService,
			"recordStorage":  imageService,
		}),
		StatisticsService: status.NewStatisticsService(status.Providers{
			"archiveStorage": archiveService,
		}),
	}

	log.Info("server startup completed")
	if err := srv.Run(); err != nil {
		log.WithError(err).Error("error starting http server")
	}
}

func ensureTmpDirectory() error {
	if _, err := os.Stat(os.TempDir()); os.IsNotExist(err) {
		if err := os.Mkdir(os.TempDir(), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
