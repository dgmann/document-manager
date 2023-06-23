package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.golang/paho"
	"github.com/otiai10/gosseract/v2"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type OCRRequest struct {
	RecordId string `json:"recordId"`
	Force    bool   `json:"force"`
}

const OCRRequestTopic = "ocrrequests"

func main() {
	mqttConString := os.Getenv("MQTT_BROKER")
	apiUrl := os.Getenv("API_URL")
	port := os.Getenv("HTTP_PORT")
	if len(port) == 0 {
		port = "8080"
	}

	log.Printf("Using API URL %s\n", apiUrl)
	log.Printf("Connecting to MQTT Broker at %s\n", mqttConString)
	conn, err := net.Dial("tcp", mqttConString)
	if err != nil {
		log.Fatalf("error opening connection to %s: %s\n", mqttConString, err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	ctx, cancel := context.WithCancel(context.Background())

	topic := "records/+"
	mqttClient := NewMQTTSubscriber(conn, "ocr-service")

	go RunHTTPServer(ctx, port, apiUrl, mqttClient)

	mqttClient.Router().RegisterHandler(topic, func(publish *mqtt.Publish) {
		var e map[string]any
		if err := json.Unmarshal(publish.Payload, &e); err != nil {
			log.Println(err)
			return
		}
		if e["type"] == "Deleted" {
			return
		}

		if err := mqttClient.Publish(ctx, OCRRequestTopic, OCRRequest{RecordId: e["id"].(string)}); err != nil {
			log.Printf("error publishing ocr request to MQTT topic: %s", err)
			return
		}
	})

	// Already spawns go routine
	mqttClient.Router().RegisterHandler(OCRRequestTopic, func(publish *mqtt.Publish) {
		var request OCRRequest
		if err := json.Unmarshal(publish.Payload, &request); err != nil {
			log.Println(err)
			return
		}

		recordUrl, err := url.JoinPath(apiUrl, "records", request.RecordId)
		if err != nil {
			log.Printf("error creating record request url: %s", err)
			return
		}
		record, err := func() (data map[string]any, err error) {
			recordResp, err := http.Get(recordUrl)
			if err != nil {
				return nil, fmt.Errorf("error fetching record: %s", err)
			}
			recordBody, err := io.ReadAll(recordResp.Body)
			defer func(Body io.ReadCloser) {
				closeErr := Body.Close()
				if closeErr != nil && err == nil {
					err = closeErr
				}
			}(recordResp.Body)
			var record map[string]any
			if err := json.Unmarshal(recordBody, &record); err != nil {
				return nil, fmt.Errorf("error parsing record json: %s", err)
			}
			return record, nil
		}()
		if err != nil {
			log.Println(err)
			return
		}

		pages, ok := record["pages"].([]interface{})
		if !ok {
			log.Println("could not extract pages")
			return
		}
		client := gosseract.NewClient()
		defer func(client *gosseract.Client) {
			err := client.Close()
			if err != nil {

			}
		}(client)
		if err := client.SetLanguage("deu", "eng"); err != nil {
			log.Fatalln(err)
		}

		// tracks whether anything was changed and an update is required
		needsUpdate := false
		updatedPages := make([]map[string]any, len(pages))
		for i, p := range pages {
			page, ok := p.(map[string]any)
			if !ok {
				log.Println("could not extract page")
				return
			}
			// If page content is already filled we do not need to scan it again
			// Force overrides this
			if content, ok := page["content"].(string); (ok && len(content) > 0) || request.Force {
				updatedPages[i] = map[string]any{
					"id": page["id"],
				}
				continue
			}

			pageUrl, ok := page["url"].(string)
			if !ok {
				log.Println("could not extract url")
				return
			}
			resp, err := http.Get(pageUrl)
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
			// If text was empty we won't update it
			if len(text) == 0 {
				updatedPages[i] = map[string]any{
					"id": page["id"],
				}
				continue
			}

			updatedPages[i] = map[string]any{
				"id":      page["id"],
				"content": text,
			}
			needsUpdate = true
		}
		if !needsUpdate {
			log.Printf("skipping record %s as no page content was changed", request.RecordId)
			return
		}
		updateUrl := recordUrl + "/pages"
		var b bytes.Buffer
		if err := json.NewEncoder(&b).Encode(updatedPages); err != nil {
			log.Printf("error encoding page update request: %s", err)
		}
		updateResp, err := http.Post(updateUrl, "application/json", &b)
		if err != nil {
			log.Printf("error updating pages at %s: %s\n", updateUrl, err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(updateResp.Body)
		if updateResp.StatusCode >= 400 {
			buf := new(strings.Builder)
			_, _ = io.Copy(buf, updateResp.Body)
			log.Printf("error updating pages. Status Code: %d, error: %s\n", updateResp.StatusCode, buf.String())
			return
		}
		log.Printf("updated pages of record %s\n", request.RecordId)

	})

	if _, err := mqttClient.Connect(ctx); err != nil {
		log.Fatalf("error connecting to MQTT broker: %s\n", err)
	}
	if err := mqttClient.Subscribe(ctx, OCRRequestTopic); err != nil {
		log.Fatalf("error subscribing to topic %s: %s\n", topic, err)
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
