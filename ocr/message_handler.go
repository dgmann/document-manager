package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/dgmann/document-manager/api/pkg/client"
	mqtt "github.com/eclipse/paho.golang/paho"
	"io"
	"log"
	"net/http"
	"ocr/internal/ocr"
	"regexp"
	"strings"
)

type OCRRequest struct {
	RecordId string `json:"recordId"`
	Force    bool   `json:"force"`
}

type CategorizationRequest struct {
	Record *api.Record `json:"record"`
}

type Handler struct {
	OCRClient ocr.Client
	ApiClient *client.HTTPClient
}

func (h *Handler) Close() error {
	return h.OCRClient.Close()
}

func backendEventHandler(ocrRequestChan chan<- OCRRequest, categorizationChan chan<- CategorizationRequest) mqtt.MessageHandler {
	return func(publish *mqtt.Publish) {
		var e api.Event[*api.Record]
		if err := json.Unmarshal(publish.Payload, &e); err != nil {
			log.Println(err)
			return
		}
		if e.Type == api.EventTypeDeleted {
			return
		}
		record := e.Data

		// Go through all pages and if any pages does not have a content, issue a OCRRequest
		for _, p := range record.Pages {
			if p.Content == nil {
				ocrRequestChan <- OCRRequest{
					RecordId: record.Id,
					Force:    false,
				}
				return
			}
		}

		if record.Category == nil || len(*record.Category) == 0 {
			categorizationChan <- CategorizationRequest{Record: record}
			return
		}
	}
}

func (h *Handler) OCRRequestHandler() mqtt.MessageHandler {
	return func(publish *mqtt.Publish) {
		var request OCRRequest
		if err := json.Unmarshal(publish.Payload, &request); err != nil {
			log.Println(err)
			return
		}

		record, err := h.ApiClient.Records.Get(request.RecordId)
		if err != nil {
			log.Println(err)
			return
		}

		pagesToProcess := make([]ocr.PageWithContent, len(record.Pages))
		for i, page := range record.Pages {
			pagesToProcess[i] = ocr.PageWithContent{Id: page.Id, Image: []byte{}}
			// If page content is already filled we do not need to scan it again
			// Force overrides this
			if page.Content != nil && !request.Force {
				continue
			}

			resp, err := http.Get(page.Url)
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
			pagesToProcess[i].Image = img
		}

		// tracks whether anything was changed and an update is required
		pagesToUpdate, err := h.OCRClient.CheckOrientation(pagesToProcess)
		if err != nil {
			log.Println(err)
			return
		}
		// if the pages need to be updated, update them first
		if len(pagesToUpdate) != 0 {
			if _, err := h.ApiClient.Records.UpdatePages(request.RecordId, pagesToUpdate); err != nil {
				log.Println(err)
			}
			log.Printf("updated pages of record %s\n", request.RecordId)
			return
		}

		pagesToUpdate, err = h.OCRClient.ExtractText(pagesToProcess)
		if err != nil {
			log.Println(err)
			return
		}
		if len(pagesToUpdate) != 0 {
			if _, err := h.ApiClient.Records.UpdatePages(request.RecordId, pagesToUpdate); err != nil {
				log.Println(err)
			}
			log.Printf("updated pages of record %s\n", request.RecordId)
			return
		}
		log.Printf("skipping record %s as no page content was changed", request.RecordId)
	}
}

func (h *Handler) CategorizationRequestHandler() mqtt.MessageHandler {
	return func(publish *mqtt.Publish) {
		var request CategorizationRequest
		if err := json.Unmarshal(publish.Payload, &request); err != nil {
			log.Println(err)
			return
		}
		contents := make([]string, len(request.Record.Pages))
		for i, page := range request.Record.Pages {
			contents[i] = *page.Content
		}
		textToSearch := strings.Join(contents, "\n")

		categories, err := h.ApiClient.Categories.All()
		if err != nil {
			log.Println(err)
			return
		}
		for _, category := range categories {
			if match(textToSearch, category.Match) {
				request.Record.Category = &category.Id
			}
		}
	}
}

func match(content string, matchConfig api.MatchConfig) bool {
	switch matchConfig.Type {
	case api.MatchTypeExact:
		matched, err := regexp.MatchString(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(matchConfig.Expression)), content)
		if err != nil {
			log.Println(err)
		}
		return matched
	case api.MatchTypeRegex:
		matched, err := regexp.MatchString(matchConfig.Expression, content)
		if err != nil {
			log.Println(err)
		}
		return matched
	case api.MatchTypeAll:
		parts := splitStringEx(matchConfig.Expression)
		for _, part := range parts {
			withoutQuotes := strings.ReplaceAll(part, `"`, "")
			matched, err := regexp.MatchString(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(withoutQuotes)), content)
			if err != nil {
				log.Println(err)
			}
			if !matched {
				return false
			}
		}
		return true
	case api.MatchTypeAny:
		parts := splitStringEx(matchConfig.Expression)
		for _, part := range parts {
			withoutQuotes := strings.ReplaceAll(part, `"`, "")
			matched, err := regexp.MatchString(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(withoutQuotes)), content)
			if err != nil {
				log.Println(err)
			}
			if matched {
				return true
			}
		}
		return false
	}
	return false
}

// splitStringEx splits a string on whitespace but keeps words grouped inside quotes
func splitStringEx(content string) []string {
	quoted := false
	return strings.FieldsFunc(content, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == ' '
	})
}
