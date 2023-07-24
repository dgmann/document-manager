package main

import (
	"context"
	"github.com/dgmann/document-manager/api/pkg/client"
	"github.com/dgmann/document-manager/pkg/opentelemetry"
	"github.com/eclipse/paho.golang/paho"
	log "github.com/sirupsen/logrus"
	"ocr/internal/ocr/tesseract"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	RecordsTopic               = "records/+"
	OCRRequestTopic            = "ocrrequests"
	CategorizationRequestTopic = "categorizationrequests"
)

func main() {
	config, err := getConfig()
	if err != nil {
		log.Fatalln(config)
	}
	log.Printf("Using API URL %s", config.ApiUrl)
	log.Printf("Connecting to MQTT Broker at %s", config.Broker)

	ctx, cancel := context.WithCancel(context.Background())

	otlProvider, err := opentelemetry.NewProvider(ctx, "ocr", config.OtelCollectorUrl)
	if err != nil {
		log.WithError(err).Warnln("error creating OpenTelemetry exporter")
	}

	ocrRequestPublishChan := make(chan OCRRequest)
	defer close(ocrRequestPublishChan)
	categorizationRequestChan := make(chan CategorizationRequest)
	defer close(categorizationRequestChan)

	apiClient, err := client.NewHTTPClient(config.ApiUrl, 3*time.Second)
	if err != nil {
		log.Fatalf("error creating API Client: %s", err)
	}
	handler := &Handler{OCRClient: tesseract.NewClient(), ApiClient: apiClient}
	defer func(h *Handler) {
		if err := h.Close(); err != nil {
			log.Printf("error closing tesseract: %s\n", err)
		}
	}(handler)
	mqttClient := NewMQTTClient(ctx, config.Broker, config.ClientId, []Subscription{
		{Topic: RecordsTopic, SubscribeOptions: paho.SubscribeOptions{QoS: 1}, Handler: backendEventHandler(ocrRequestPublishChan, categorizationRequestChan)},
		{Topic: OCRRequestTopic, SubscribeOptions: paho.SubscribeOptions{QoS: 1}, Handler: handler.OCRRequestHandler()},
		{Topic: CategorizationRequestTopic, SubscribeOptions: paho.SubscribeOptions{QoS: 1}, Handler: handler.CategorizationRequestHandler()},
	})
	if err := mqttClient.Connect(ctx); err != nil {
		log.Fatalf("error connecting subscriber: %s", err)
	}

	go RunHTTPServer(ctx, config, ocrRequestPublishChan)
	go func() {
		err := Publish(ctx, mqttClient, OCRRequestTopic, ocrRequestPublishChan)
		if err != nil {
			log.Fatalf("OCRRequest publisher error: %s", err)
		}
	}()
	go func() {
		err := Publish(ctx, mqttClient, CategorizationRequestTopic, categorizationRequestChan)
		if err != nil {
			log.Fatalf("CategorizationRequest publisher error: %s", err)
		}
	}()

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
			_ = mqttClient.Disconnect(ctx)
			if otlProvider != nil {
				_ = otlProvider.Shutdown(ctx)
			}
		}()
		cancel()
	}()
	<-ctx.Done()
}
