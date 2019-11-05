package http

import (
	"fmt"
	"github.com/dgmann/document-manager/api/datastore"
	"github.com/dgmann/document-manager/api/datastore/mock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRecordController_All(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	mockRecordRepository := new(mock.RecordService)
	controller := RecordController{records: mockRecordRepository}
	record := datastore.NewRecord(datastore.CreateRecord{Sender: "mock", Status: datastore.StatusInbox})
	mockRecordRepository.On("Query", req.Context(), datastore.NewRecordQuery(), datastore.NewQueryOptions()).Return([]datastore.Record{*record}, nil)

	controller.All(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, fmt.Sprintf(`[%s]`, buildJSONResponseForRecord(*record)), w.Body.String())
}

func buildJSONResponseForRecord(record datastore.Record) string {
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
	}`, record.Id.Hex(), record.ReceivedAt.Format(time.RFC3339), record.UpdatedAt.Format(time.RFC3339Nano), datastore.StatusInbox)
}
