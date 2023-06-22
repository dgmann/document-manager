package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ScanJob struct {
	Full bool `json:"full,omitempty"`
}

type RecordResponse struct {
	Id string `json:"id"`
}

func RunHTTPServer(ctx context.Context, port string, apiUrl string, mqttClient *MQTTClient) {
	srv := &http.Server{Addr: ":" + port}
	recordUrl, err := url.JoinPath(apiUrl, "records")
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("OCR Service")); err != nil {
			log.Printf(err.Error())
		}
	})

	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		response := json.NewEncoder(w)

		respond := func(statusCode int, body interface{}) {
			w.WriteHeader(statusCode)
			_ = response.Encode(body)
		}

		if r.Method != http.MethodPost {
			respond(http.StatusMethodNotAllowed, map[string]any{"error": "use POST instead"})
			return
		}

		var job ScanJob
		if err := json.NewDecoder(r.Body).Decode(&job); err != nil && err != io.EOF {
			respond(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}

		recordRequest, err := http.Get(recordUrl)
		if err != nil {
			respond(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}

		var recordResponse []RecordResponse
		if err := json.NewDecoder(recordRequest.Body).Decode(&recordResponse); err != nil {
			respond(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}

		for _, record := range recordResponse {
			if err := mqttClient.Publish(ctx, OCRRequestTopic, OCRRequest{RecordId: record.Id, Force: job.Full}); err != nil {
				return
			}
		}

		respond(http.StatusAccepted, map[string]any{"message": fmt.Sprintf("published %d ocr requests", len(recordResponse))})
		return
	})
	go func() {
		select {
		case <-ctx.Done():
			if err := srv.Shutdown(context.Background()); err != nil {
				log.Println(err)
			}
			return
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
