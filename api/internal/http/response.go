package http

import (
	"encoding/json"
	"net/http"
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

// Generic HTTP Response
// swagger:ignore
type Response struct {
	StatusCode int
	writer     http.ResponseWriter `json:"-"`
}

// Data HTTP Response
//
// A DataResponse is a Response which provides additional data as a payload.
//
// swagger:response
type DataResponse struct {
	Response
	Data interface{}
}

func (r *DataResponse) WriteJSON() {
	writeJSON(r.writer, r.Data, r.StatusCode)
}

// Binary HTTP Response
//
// A BinaryResponse is a Response which returns binary data.
//
// swagger:response
type BinaryResponse struct {
	Response
	Data []byte
}

func NewBinaryResponseWithStatus(writer http.ResponseWriter, data []byte, code int) *BinaryResponse {
	response := Response{writer: writer, StatusCode: code}
	return &BinaryResponse{Response: response, Data: data}
}

func (r *BinaryResponse) Write() {
	if _, err := r.writer.Write(r.Data); err != nil {
		r.writer.WriteHeader(500)
		return
	}
	r.writer.WriteHeader(r.StatusCode)
}

// Error HTTP Response
//
// A ErrorResponse indicates an error and provides the error message.
//
// swagger:response
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
