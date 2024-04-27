package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"time"
)

func NewResponse(writer http.ResponseWriter, data interface{}) *DataResponse {
	return NewResponseWithStatus(writer, data, http.StatusOK)
}

func NewResponseWithStatus(writer http.ResponseWriter, data interface{}, code int) *DataResponse {
	response := Response{writer: writer, StatusCode: code}
	return &DataResponse{Response: response, Data: data}
}

func NewErrorResponse(writer http.ResponseWriter, err error, code int) *ErrorResponse {
	response := Response{writer: writer, StatusCode: code}
	return &ErrorResponse{Response: response, Error: err.Error()}
}

// Response serves as the base and most generic form of response
type Response struct {
	StatusCode int
	writer     http.ResponseWriter
}

func (r *Response) SetEtag(etag ETag) *Response {
	r.writer.Header().Set(etag.Header())
	return r
}

func (r *Response) Write() {
	r.writer.WriteHeader(r.StatusCode)
}

// DataResponse extends Response a generic payload.
type DataResponse struct {
	Response
	Data interface{}
}

func (r *DataResponse) WriteJSON() {
	writeJSON(r.writer, r.Data, r.StatusCode)
}

func (r *DataResponse) SetEtag(etag ETag) *DataResponse {
	r.Response.SetEtag(etag)
	return r
}

// BinaryResponse extends Response with a binary data payload.
type BinaryResponse struct {
	Response
	Data []byte
}

func NewBinaryResponseWithStatus(writer http.ResponseWriter, data []byte, code int) *BinaryResponse {
	response := Response{writer: writer, StatusCode: code}
	return &BinaryResponse{Response: response, Data: data}
}

func (r *BinaryResponse) SetContentTypeFromExtension(ex string) *BinaryResponse {
	mimeType := mime.TypeByExtension(ex)
	r.writer.Header().Set("Content-Type", mimeType)
	return r
}

type ETag struct {
	time time.Time
}

func NewEtag(t time.Time) ETag {
	return ETag{time: t}
}

func (e ETag) String() string {
	return strconv.FormatInt(e.time.UTC().UnixNano(), 10)
}

func (e ETag) Header() (string, string) {
	return "ETag", e.String()
}

var ErrNoEtagHeader = errors.New("no If-Match header found")

func EtagFromHeader(header http.Header) (time.Time, error) {
	value := header.Get("If-Match")
	if len(value) == 0 {
		return time.Time{}, ErrNoEtagHeader
	}
	timestamp, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing If-Match header: %w", err)
	}
	return time.Unix(0, timestamp), nil
}

func (r *BinaryResponse) Write() {
	if _, err := r.writer.Write(r.Data); err != nil {
		r.writer.WriteHeader(500)
		return
	}
	r.writer.WriteHeader(r.StatusCode)
}

// ErrorResponse indicates an error and provides the error message.
type ErrorResponse struct {
	Response
	Error string
}

func (r *ErrorResponse) WriteJSON() {
	message := map[string]string{"error": r.Error}
	writeJSON(r.writer, message, r.StatusCode)
}

func writeJSON(w http.ResponseWriter, message interface{}, code int) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
