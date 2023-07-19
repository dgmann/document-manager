package client

import (
	"bytes"
	"encoding/json"
	"github.com/dgmann/document-manager/api/pkg/api"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpUploader_Upload(t *testing.T) {
	sender := "Test"
	receivedAt := time.Now()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, sender, req.FormValue("sender"))
		assert.Equal(t, receivedAt.Format(time.RFC3339), req.FormValue("receivedAt"))
		assert.Equal(t, req.URL.String(), "/records")
		// Send response to be tested
		rw.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(rw).Encode(api.Record{})
	}))
	defer server.Close()

	record := NewRecord{File: bytes.NewBufferString(""), ReceivedAt: receivedAt, Sender: sender, Date: &time.Time{}}
	client, err := NewHTTPClient(server.URL, 10*time.Second)
	if err != nil {
		t.Error(err)
	}
	_, err = client.Records.Create(&record)
	assert.NoError(t, err)
}

func TestHttpUploader_Upload_Failed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	record := NewRecord{File: bytes.NewBufferString(""), ReceivedAt: time.Now(), Sender: "Test", Date: &time.Time{}}
	client, err := NewHTTPClient(server.URL, 10*time.Second)
	if err != nil {
		t.Error(err)
	}
	_, err = client.Records.Create(&record)
	assert.Error(t, err)
}
