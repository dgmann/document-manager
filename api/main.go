package main

import (
	"context"
	"github.com/dgmann/document-manager/api/internal/datastore/mongo"
	"github.com/dgmann/document-manager/api/internal/event"
	"github.com/dgmann/document-manager/api/internal/http"
	"github.com/dgmann/document-manager/api/internal/pdf/grpc"
	"github.com/dgmann/document-manager/api/internal/status"
	"github.com/dgmann/document-manager/api/internal/storage/filesystem"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/dgmann/document-manager/pkg/opentelemetry"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
	)))
}

const ServiceName = "backend"

func main() {
	config := ConfigFromEnv()
	ctx, cancel := context.WithCancel(context.Background())

	otlProvider, err := opentelemetry.NewProvider(ctx, ServiceName, config.OtelCollectorUrl)
	if err != nil {
		log.WithContext(ctx).WithError(err).Warnln("error creating OpenTelemetry exporter")
	}
	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		log.WithContext(ctx).WithError(err).Warnln("error initializing runtime metrics")
	}

	if err := ensureTmpDirectory(); err != nil {
		log.WithContext(ctx).WithError(err).Error("error while creating tmp directory")
		return
	}
	log.WithContext(ctx).
		WithFields(log.Fields{"host": config.Database.Host, "port": config.Database.Port, "database": config.Database.Name}).
		Info("connecting to database")

	client := func() *mongo.Client {
		client := mongo.NewClient(config.Database)
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := client.Connect(ctx); err != nil {
			log.WithContext(ctx).WithError(err).Fatalln("database cannot be reached")
		}
		return client
	}()

	if err := client.CreateIndexes(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("error setting indices")
	}

	imageService, err := filesystem.NewImageService(config.RecordDirectory)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error creating image service")
		return
	}
	archiveService, err := filesystem.NewArchiveService(config.ArchiveDirectory)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error creating archive service")
		return
	}
	categoryService := mongo.NewCategoryService(mongo.NewCollection(client.Categories()), mongo.NewCollection(client.Records()))

	pdfProcessor, err := grpc.NewPDFProcessor(config.PdfProcessorUrl, imageService, categoryService)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error connecting to pdf processor service")
		return
	}

	websocketService := event.NewWebsocketEventService[*api.Record]()
	mqttBrokerUrl, err := url.Parse(config.MQTTBroker)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("error opening connection to %s\n", config.MQTTBroker)
		return
	}
	mqttService := func() *event.MQTTEventSender[*api.Record] {
		mqttService := event.NewMQTTEventSender[*api.Record](mqttBrokerUrl, config.MQTTClientId)
		if err := mqttService.Connect(ctx); err != nil {
			log.WithContext(ctx).WithError(err).Fatalln("error connecting to MQTT Broker")
		}
		return mqttService
	}()
	eventService := event.NewMultiEventSender[*api.Record](websocketService, mqttService)

	tagService := mongo.NewTagService(client.Records())

	recordService := mongo.NewRecordService(mongo.RecordServiceConfig{
		Records: mongo.NewCollection(client.Records()),
		Events:  eventService,
		Images:  imageService,
		Pdfs:    archiveService,
	})
	healthService := status.NewHealthService(status.HealthCheckables{
		"database":       client,
		"pdfProcessor":   pdfProcessor,
		"archiveStorage": archiveService,
		"recordStorage":  imageService,
	})
	statisticsService := status.NewStatisticsService(status.Providers{
		"archiveStorage": archiveService,
	})
	srv := http.NewServer(
		ServiceName,
		http.WithRecordController(recordService, imageService, archiveService, pdfProcessor),
		http.WithPatientController(recordService, imageService, categoryService, tagService),
		http.WithCategoryController(categoryService),
		http.WithTagController(tagService),
		http.WithArchiveController(archiveService),
		http.WithHealthController(healthService),
		http.WithStatisticController(statisticsService),
		http.WithExportController(recordService, pdfProcessor),
		http.WithNotificationController(websocketService),
	)

	go func() {
		if err := srv.Run(config.Port); err != nil {
			log.WithContext(ctx).WithError(err).Error("error starting http server")
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
			if err := srv.Shutdown(ctx); err != nil {
				log.Printf("error shutting down http server: %s", err)
			}
			if err := mqttService.Disconnect(ctx); err != nil {
				log.Printf("error disconnecting MQTT: %s", err)
			}
			if otlProvider != nil {
				if err := otlProvider.Shutdown(ctx); err != nil {
					log.Printf("error shutting down tracer provider: %s", err)
				}
			}
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
