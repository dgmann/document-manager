package http

import (
	"fmt"
	"github.com/dgmann/document-manager/api/internal/datastore"
	"github.com/dgmann/document-manager/api/internal/datastore/mock"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRecordController_All(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	mockRecordRepository := mock.NewRecordService()
	controller := RecordController{records: mockRecordRepository}
	record := api.NewRecord(api.CreateRecord{Sender: "mock", Status: api.StatusInbox})
	mockRecordRepository.On("Query", req.Context(), datastore.NewRecordQuery(), datastore.NewQueryOptions()).Return([]api.Record{*record}, nil)

	controller.All(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, fmt.Sprintf(`[%s]`, buildJSONResponseForRecord(*record)), w.Body.String())
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
