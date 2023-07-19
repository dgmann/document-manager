package http

import (
	"encoding/json"
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

func (r *BinaryResponse) SetEtag(etag string) *BinaryResponse {
	r.writer.Header().Set("ETag", etag)
	return r
}

func ETag(t time.Time) string {
	return strconv.FormatInt(t.UTC().Unix(), 10)
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
