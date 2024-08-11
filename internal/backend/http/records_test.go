package http

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgmann/document-manager/internal/backend/datastore"
	"github.com/dgmann/document-manager/internal/backend/datastore/mock"
	"github.com/dgmann/document-manager/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestRecordController_All(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	mockRecordRepository := mock.NewRecordService()
	controller := RecordController{records: mockRecordRepository}
	status := api.StatusInbox
	record := api.Record{Sender: "mock", Status: &status}
	mockRecordRepository.On("Query", req.Context(), datastore.NewRecordQuery(), datastore.NewQueryOptions()).Return([]api.Record{record}, nil)

	controller.All(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, fmt.Sprintf(`[%s]`, buildJSONResponseForRecord(record)), w.Body.String())
}

func buildJSONResponseForRecord(record api.Record) string {
	return fmt.Sprintf(`{
		"id":"%[1]s",
		"archivedPDF": "http://example.com/api/archive/%[1]s", 
		"category": null, 
		"comment": "", 
		"date": null,
		"pages": [], 
		"patientId": "", 
		"receivedAt": "%s",
		"updatedAt": "%s",
		"sender": "mock", 
		"status": "%s", 
		"tags": []
	}`, record.Id, record.ReceivedAt.Format(time.RFC3339), record.UpdatedAt.Format(time.RFC3339Nano), api.StatusInbox)
}
