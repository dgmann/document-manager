package main

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/api/internal/datastore/mongo"
	"github.com/dgmann/document-manager/api/internal/event"
	"github.com/dgmann/document-manager/api/internal/http"
	"github.com/dgmann/document-manager/api/internal/pdf/grpc"
	"github.com/dgmann/document-manager/api/internal/status"
	"github.com/dgmann/document-manager/api/internal/storage/filesystem"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
	"os/signal"
	"syscall"
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
	log.WithFields(log.Fields{"host": config.Database.Host, "port": config.Database.Port, "database": config.Database.Name}).Info("connecting to database")
	ctx, cancel := context.WithCancel(context.Background())

	client := func() *mongo.Client {
		client := mongo.NewClient(config.Database)
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := client.Connect(ctx); err != nil {
			log.WithError(err).Error("database cannot be reached")
		}
		return client
	}()

	if err := client.CreateIndexes(ctx); err != nil {
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

	websocketService := event.NewWebsocketEventService()

	mqttBrokerUrl, err := url.Parse(config.MQTTBroker)
	if err != nil {
		log.WithError(err).Fatalf("error opening connection to %s\n", config.MQTTBroker)
	}
	mqttService := func() *event.MQTTEventSender {
		mqttService := event.NewMQTTEventSender(mqttBrokerUrl, config.MQTTClientId)
		if err := mqttService.Connect(ctx); err != nil {
			log.WithError(err).Fatalln("error connecting to MQTT Broker")
		}
		return mqttService
	}()

	eventService := event.NewMultiEventSender(websocketService, mqttService)

	tagService := mongo.NewTagService(client.Records())

	srv := http.Server{
		Port:            config.Port,
		EventService:    websocketService,
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

	go func() {
		if err := srv.Run(); err != nil {
			log.WithError(err).Error("error starting http server")
		}
	}()
	log.Info("server startup completed")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)
	go func() {
		<-signalChan
		log.Print("os.Interrupt - shutting down...\n")
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_ = srv.Shutdown(ctx)
			_ = mqttService.Disconnect(ctx)
		}()
		cancel()
	}()
	<-ctx.Done()
}

func ensureTmpDirectory() error {
	if _, err := os.Stat(os.TempDir()); os.IsNotExist(err) {
		if err := os.Mkdir(os.TempDir(), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
