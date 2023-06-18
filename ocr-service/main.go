package main

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.golang/paho"
	"github.com/otiai10/gosseract/v2"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	mqttConString := os.Getenv("MQTT_BROKER")

	conn, err := net.Dial("tcp", mqttConString)
	if err != nil {
		log.Fatalf("error opening connection: %s\n", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	client := gosseract.NewClient()
	defer func(client *gosseract.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)
	if err := client.SetLanguage("deu", "eng"); err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	topic := "records/+"
	mqttClient := NewMQTTSubscriber(conn, "ocr-service")
	// Already spawns go routine
	mqttClient.Router().RegisterHandler(topic, func(publish *mqtt.Publish) {
		var e map[string]any
		if err := json.Unmarshal(publish.Payload, &e); err != nil {
			log.Println(err)
			return
		}
		if e["type"] == "Deleted" {
			return
		}
		record, ok := e["data"].(map[string]any)
		if !ok {
			log.Println("could not extract record")
			return
		}
		pages, ok := record["pages"].([]interface{})
		if !ok {
			log.Println("could not extract pages")
			return
		}
		for _, p := range pages {
			page, ok := p.(map[string]any)
			if !ok {
				log.Println("could not extract page")
				return
			}
			url, ok := page["url"].(string)
			if !ok {
				log.Println("could not extract url")
				return
			}
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
				return
			}
			img, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err != nil {
				log.Println(err)
				return
			}
			if err != client.SetImageFromBytes(img) {
				log.Println(err)
			}
			text, err := client.Text()
			if err != nil {
				log.Println(err)
			}
			fmt.Println(text)
		}
	})

	if _, err := mqttClient.Connect(ctx); err != nil {
		log.Fatalf("error connecting to MQTT broker: %s\n", err)
	}
	if err := mqttClient.Subscribe(ctx, topic); err != nil {
		log.Fatalf("error subscribing to topic %s: %s\n", topic, err)
	}
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
		cancel()
	}()
	<-ctx.Done()
}
