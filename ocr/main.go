package main

import (
	"context"
	"github.com/eclipse/paho.golang/paho"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type OCRRequest struct {
	RecordId string `json:"recordId"`
	Force    bool   `json:"force"`
}

const (
	RecordsTopic    = "records/+"
	OCRRequestTopic = "ocrrequests"
)

func main() {
	config, err := getConfig()
	if err != nil {
		log.Fatalln(config)
	}
	log.Printf("Using API URL %s\n", config.ApiUrl)
	log.Printf("Connecting to MQTT Broker at %s\n", config.Broker)

	ctx, cancel := context.WithCancel(context.Background())

	ocrRequestPublishChan := make(chan OCRRequest)
	defer close(ocrRequestPublishChan)

	client := NewMQTTClient(ctx, config.Broker, config.ClientId, []Subscription{
		{Topic: RecordsTopic, SubscribeOptions: paho.SubscribeOptions{QoS: 1}, Handler: handleBackendEvent(ocrRequestPublishChan)},
		{Topic: OCRRequestTopic, SubscribeOptions: paho.SubscribeOptions{QoS: 1}, Handler: handlerOCRRequest(config.ApiUrl)},
	})
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("error connecting subscriber: %s", err)
	}

	go RunHTTPServer(ctx, config, ocrRequestPublishChan)
	go func() {
		err := client.Run(ctx, OCRRequestTopic, ocrRequestPublishChan)
		if err != nil {
			log.Fatalf("publisher error: %s", err)
		}
	}()

	log.Println("listening...")

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
			_ = client.Disconnect(ctx)
		}()
		cancel()
	}()
	<-ctx.Done()
}
