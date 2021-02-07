package grpc

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dgmann/document-manager/api/datastore"
	"github.com/dgmann/document-manager/api/storage"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
)

func NewDocument(title string, records []datastore.Record, images storage.ImageService, categories []datastore.Category) (*processor.Document, error) {
	recordsGroupedByCategory := make(map[string][]datastore.Record)
	categoryMap := createCategoryLookupMap(categories)
	for _, record := range records {
		category := strings.Title(string(*record.Status))
		if record.Category != nil {
			category = categoryMap[*record.Category].Name
		}
		if _, ok := recordsGroupedByCategory[category]; !ok {
			recordsGroupedByCategory[category] = make([]datastore.Record, 0)
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

			return !date1.Before(date2)
		})
		subdocuments := make([]*processor.Document, 0, len(grouped))
		for _, record := range grouped {
			title := fmt.Sprintf("Empfangen: %s", record.ReceivedAt.Format("02.01.2006 15:04"))
			if record.Date != nil {
				title = record.Date.Format("02.01.2006")
			}
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

func createCategoryLookupMap(categories []datastore.Category) map[string]datastore.Category {
	res := make(map[string]datastore.Category)
	for _, cat := range categories {
		res[cat.Id] = cat
	}
	return res
}
