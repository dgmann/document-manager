package http

import (
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"github.com/dgmann/document-manager/api/app/mock"
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
	record := app.NewRecord(app.CreateRecord{Sender: "mock", Status: app.StatusInbox})
	mockRecordRepository.On("Query", mock.Anything, map[string]interface{}{}).Return([]app.Record{*record}, nil)

	controller.All(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, fmt.Sprintf(`[%s]`, buildJSONResponseForRecord(*record)), w.Body.String())
}

func buildJSONResponseForRecord(record app.Record) string {
	return fmt.Sprintf(`{
		"id":"%[1]s",
		"archivedPDF": "http://example.com/archive/%[1]s", 
		"category": null, 
		"comment": "", 
		"date": null,
		"pages": [], 
		"patientId": "", 
		"receivedAt": "%s", 
		"sender": "mock", 
		"status": "%s", 
		"tags": []
	}`, record.Id.Hex(), record.ReceivedAt.Format(time.RFC3339), app.StatusInbox)
}
