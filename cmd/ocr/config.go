package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	Broker           *url.URL
	ApiUrl           string
	HttpPort         string
	ClientId         string
	OtelCollectorUrl string
}

func getConfig() (Config, error) {
	envBroker := os.Getenv("MQTT_BROKER")
	mqttConString, err := url.Parse(envBroker)
	if err != nil {
		return Config{}, fmt.Errorf("environmental variable %s must be a valid URL (%w)", envBroker, err)
	}
	apiUrl := os.Getenv("API_URL")
	port := os.Getenv("HTTP_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	clientId := os.Getenv("MQTT_CLIENT_ID")
	if len(clientId) == 0 {
		clientId = "ocr-service-" + strconv.Itoa(rand.Int())
	}
	return Config{
		Broker:           mqttConString,
		ApiUrl:           apiUrl,
		HttpPort:         port,
		ClientId:         clientId,
		OtelCollectorUrl: os.Getenv("OTEL_COLLECTOR_URL"),
	}, nil
}
