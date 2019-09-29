package grpc

import (
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"sort"
	"time"
)

func NewDocument(title string, records []app.Record, images app.ImageService) (*processor.Document, error) {
	recordsGroupedByCategory := make(map[string][]app.Record)
	for _, record := range records {
		category := ""
		if record.Category != nil {
			category = *record.Category
		}
		if _, ok := recordsGroupedByCategory[category]; !ok {
			recordsGroupedByCategory[category] = make([]app.Record, 0)
		}
		recordsGroupedByCategory[category] = append(recordsGroupedByCategory[category], record)
	}

	documents := make([]*processor.Document, 0, len(recordsGroupedByCategory))
	for category, grouped := range recordsGroupedByCategory {
		sort.Slice(grouped, func(i, j int) bool {
			var date1, date2 time.Time
			if grouped[i].Date != nil {
				date1 = *grouped[i].Date
			} else {
				date1 = time.Time{}
			}

			if grouped[j].Date != nil {
				date2 = *grouped[j].Date
			} else {
				date2 = time.Time{}
			}

			return date1.Before(date2)
		})
		subdocuments := make([]*processor.Document, 0, len(grouped))
		for _, record := range grouped {
			title := record.Id.Hex()
			imagesForRecord, err := images.Get(record.Id.Hex())
			if err != nil {
				return nil, fmt.Errorf("error fetching images for record %s: %w", record.Id.Hex(), err)
			}
			pages := make([]*processor.Image, len(record.Pages))
			for i, page := range record.Pages {
				if content, ok := imagesForRecord[page.Id]; ok {
					pages[i] = &processor.Image{Format: content.Format, Content: content.Image}
				} else {
					return nil, fmt.Errorf("image %s of reocrd %s could not be found: %w", page.Id, record.Id.Hex(), err)
				}
			}
			entry := &processor.Document{Title: title, Pages: pages}
			subdocuments = append(subdocuments, entry)
		}
		doc := &processor.Document{Title: category, Documents: subdocuments}
		documents = append(documents, doc)
	}

	return &processor.Document{
		Title:     title,
		Documents: documents,
		Pages:     make([]*processor.Image, 0),
	}, nil
}
