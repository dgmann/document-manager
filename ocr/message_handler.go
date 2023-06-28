package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgmann/gosseract"
	mqtt "github.com/eclipse/paho.golang/paho"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func handleBackendEvent(publishChan chan<- OCRRequest) mqtt.MessageHandler {
	return func(publish *mqtt.Publish) {
		var e map[string]any
		if err := json.Unmarshal(publish.Payload, &e); err != nil {
			log.Println(err)
			return
		}
		if e["type"] == "Deleted" {
			return
		}

		publishChan <- OCRRequest{RecordId: e["id"].(string)}
	}
}

type PageWithContent struct {
	Id    string
	Image []byte
}

func handlerOCRRequest(apiUrl string) mqtt.MessageHandler {
	return func(publish *mqtt.Publish) {
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

		pagesToProcess := make([]PageWithContent, len(pages))

		for i, p := range pages {
			page, ok := p.(map[string]any)
			if !ok {
				log.Println("could not extract page")
				return
			}
			// If page content is already filled we do not need to scan it again
			// Force overrides this
			if content, ok := page["content"].(string); ok && len(content) > 0 && !request.Force {
				pagesToProcess[i] = PageWithContent{Id: page["id"].(string), Image: []byte{}}
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
			pagesToProcess[i] = PageWithContent{Id: page["id"].(string), Image: img}
		}

		// tracks whether anything was changed and an update is required
		updatedPages, err := checkOrientation(pagesToProcess)
		if err != nil {
			log.Println(err)
			return
		}
		if len(updatedPages) != 0 {
			if err := updatePages(recordUrl, updatedPages); err != nil {
				log.Println(err)
			}
			log.Printf("updated pages of record %s\n", request.RecordId)
			return
		}

		updatedPages, err = extractText(pagesToProcess)
		if err != nil {
			log.Println(err)
			return
		}
		if len(updatedPages) != 0 {
			if err := updatePages(recordUrl, updatedPages); err != nil {
				log.Println(err)
			}
			log.Printf("updated pages of record %s\n", request.RecordId)
			return
		}
		log.Printf("skipping record %s as no page content was changed", request.RecordId)
	}
}

func checkOrientation(pages []PageWithContent) ([]map[string]any, error) {
	client := gosseract.NewClient()
	defer func(client *gosseract.Client) {
		err := client.Close()
		if err != nil {
			log.Printf("error closing tesseract: %s\n", err)
		}
	}(client)

	needsUpdate := false
	updatedPages := make([]map[string]any, len(pages))
	for i, p := range pages {
		updatedPages[i] = map[string]any{
			"id": p.Id,
		}

		// If there is no image to process, nothing to do
		if len(p.Image) == 0 {
			continue
		}

		if err := client.SetImageFromBytes(p.Image); err != nil {
			return nil, err
		}
		if err := client.SetLanguage("osd"); err != nil {
			log.Fatalln(err)
		}
		if err := client.SetPageSegMode(gosseract.PSM_OSD_ONLY); err != nil {
			log.Fatalln()
		}
		osdResult, err := client.DetectOrientationScript()
		if err != nil {
			return nil, err
		}
		if osdResult.OrientationDegree == 0 {
			continue
		}
		log.Printf("page %s is rotated %d degree with confidence %f. Correcting orientation.", p.Id, osdResult.OrientationDegree, osdResult.OrientationConfidence)
		needsUpdate = true
		updatedPages[i] = map[string]any{
			"id":     p.Id,
			"rotate": osdResult.OrientationDegree * -1,
		}
	}
	if !needsUpdate {
		return nil, nil
	}
	return updatedPages, nil
}

func extractText(pages []PageWithContent) ([]map[string]any, error) {
	client := gosseract.NewClient()
	defer func(client *gosseract.Client) {
		err := client.Close()
		if err != nil {
			log.Printf("error closing tesseract: %s\n", err)
		}
	}(client)

	needsUpdate := false
	updatedPages := make([]map[string]any, len(pages))
	for i, p := range pages {
		updatedPages[i] = map[string]any{
			"id": p.Id,
		}

		// If there is no image to process, we just append the page without modification
		if len(p.Image) == 0 {
			continue
		}
		if err := client.SetImageFromBytes(p.Image); err != nil {
			return nil, err
		}
		if err := client.SetLanguage("deu", "eng"); err != nil {
			log.Fatalln(err)
		}
		if err := client.SetPageSegMode(gosseract.PSM_AUTO); err != nil {
			log.Fatalln(err)
		}
		text, err := client.Text()
		if err != nil {
			return nil, err
		}
		// If text was empty we won't update it
		if len(text) == 0 {
			continue
		}

		updatedPages[i] = map[string]any{
			"id":      p.Id,
			"content": text,
		}
		needsUpdate = true
	}
	if !needsUpdate {
		return nil, nil
	}
	return updatedPages, nil
}

func updatePages(recordUrl string, updatedPages []map[string]any) error {
	updateUrl := recordUrl + "/pages"
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(updatedPages); err != nil {
		return fmt.Errorf("error encoding page update request: %s", err)
	}
	updateResp, err := http.Post(updateUrl, "application/json", &b)
	if err != nil {
		return fmt.Errorf("error updating pages at %s: %s\n", updateUrl, err)
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
		return fmt.Errorf("error updating pages. Status Code: %d, error: %s\n", updateResp.StatusCode, buf.String())
	}
	return nil
}
