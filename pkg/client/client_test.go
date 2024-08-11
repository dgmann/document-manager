package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgmann/document-manager/pkg/api"

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

	record := api.NewRecord{File: bytes.NewBufferString(""), ReceivedAt: receivedAt, Sender: sender, Date: &time.Time{}}
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

	record := api.NewRecord{File: bytes.NewBufferString(""), ReceivedAt: time.Now(), Sender: "Test", Date: &time.Time{}}
	client, err := NewHTTPClient(server.URL, 10*time.Second)
	if err != nil {
		t.Error(err)
	}
	_, err = client.Records.Create(&record)
	assert.Error(t, err)
}

func TestHandleResponse(t *testing.T) {
	type args struct {
		res *http.Response
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    T
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[api.Category]{
		{"Status: 200", args{mockResponse(http.StatusOK, api.Category{Id: "1"})}, api.Category{Id: "1"}, func(t assert.TestingT, err error, i ...interface{}) bool {
			return false
		}},
		{"Status: 404", args{mockResponse(http.StatusNotFound, api.Category{})}, api.Category{}, func(t assert.TestingT, err error, i ...interface{}) bool {
			return !errors.Is(err, ErrNotFound)
		}},
		{"Status: 409", args{mockResponse(http.StatusConflict, api.Category{})}, api.Category{}, func(t assert.TestingT, err error, i ...interface{}) bool {
			return !errors.Is(err, ErrAlreadyExists)
		}},
		{"Status: Any Error", args{mockResponse(http.StatusInternalServerError, api.Category{})}, api.Category{}, func(t assert.TestingT, err error, i ...interface{}) bool {
			return !errors.Is(err, ErrGeneric)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleResponse[api.Category](tt.args.res)
			if !tt.wantErr(t, err, fmt.Sprintf("HandleResponse(%v)", tt.args.res)) {
				return
			}
			assert.Equalf(t, tt.want, got, "HandleResponse(%v)", tt.args.res)
		})
	}
}

func mockResponse(statusCode int, payload interface{}) *http.Response {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(payload)
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(body),
	}
}
