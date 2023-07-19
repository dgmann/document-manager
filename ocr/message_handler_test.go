package main

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/eclipse/paho.golang/paho"
	"testing"
)

const TestRecordId = "1"

var (
	TestContent  = "Hello World"
	TestCategory = "testCategory"
)

func Test_handleBackendEvent(t *testing.T) {
	type args struct {
		event api.Event[*api.Record]
	}
	type want struct {
		ocrRequest            bool
		categorizationRequest bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"OCR Required Single Page",
			args{createTestEvent(nil, []api.Page{{Content: nil}})},
			want{ocrRequest: true},
		},
		{"OCR Required Multi Page",
			args{createTestEvent(nil, []api.Page{{Content: nil}, {Content: &TestContent}})},
			want{ocrRequest: true},
		},
		{"Categorization Required Single Page",
			args{createTestEvent(nil, []api.Page{{Content: &TestContent}})},
			want{categorizationRequest: true},
		},
		{"Categorization Required Multi Page",
			args{createTestEvent(nil, []api.Page{{Content: &TestContent}, {Content: &TestContent}})},
			want{categorizationRequest: true},
		},
		{"Nothing Required Single Page",
			args{createTestEvent(&TestCategory, []api.Page{{Content: &TestContent}})},
			want{},
		},
		{"Nothing Required Multi Page",
			args{createTestEvent(&TestCategory, []api.Page{{Content: &TestContent}, {Content: &TestContent}})},
			want{},
		},
		{"Skip Delete Event",
			args{api.Event[*api.Record]{Type: api.EventTypeDeleted}},
			want{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ocrRequestChan := make(chan OCRRequest, 1)
			categorizationChan := make(chan CategorizationRequest, 1)
			handler := backendEventHandler(ocrRequestChan, categorizationChan)
			payload, _ := json.Marshal(tt.args.event)
			publish := paho.Publish{Payload: payload}

			handler(&publish)
			close(ocrRequestChan)
			close(categorizationChan)

			if got := chanToSlice(ocrRequestChan); (len(got) > 0) != tt.want.ocrRequest {
				t.Errorf("OCRRequests = %v, want %v", got, tt.want)
			}

			if got := chanToSlice(categorizationChan); (len(got) > 0) != tt.want.categorizationRequest {
				t.Errorf("CategorizationRequests = %v, want %v", got, tt.want)
			}
		})
	}
}

func chanToSlice[T any](elements chan T) []T {
	s := make([]T, 0, len(elements))
	for element := range elements {
		s = append(s, element)
	}
	return s
}

func createTestEvent(category *string, pages []api.Page) api.Event[*api.Record] {
	return api.Event[*api.Record]{
		Id: TestRecordId, Type: api.EventTypeCreated, Data: &api.Record{
			Id:       TestRecordId,
			Category: category,
			Pages:    pages,
		},
	}
}
