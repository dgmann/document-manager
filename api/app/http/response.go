package http

import (
	"encoding/json"
	"net/http"
)

func NewResponse(writer http.ResponseWriter, data interface{}) *Response {
	return NewResponseWithStatus(writer, data, http.StatusOK)
}

func NewResponseWithStatus(writer http.ResponseWriter, data interface{}, code int) *Response {
	return &Response{Data: data, writer: writer, StatusCode: code}
}

func NewErrorResponse(writer http.ResponseWriter, err error, code int) *Response {
	return &Response{Data: map[string]string{"error": err.Error()}, writer: writer, StatusCode: code}
}

type Response struct {
	StatusCode int
	Data       interface{}
	writer     http.ResponseWriter
}

func (r *Response) WriteJSON() {
	r.writer.WriteHeader(r.StatusCode)
	r.writer.Header().Add("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(r.writer).Encode(r.Data); err != nil {
		r.writer.WriteHeader(500)
	}
}

func (r *Response) Write() {
	r.writer.WriteHeader(r.StatusCode)
	if _, err := r.writer.Write(r.Data.([]byte)); err != nil {
		r.writer.WriteHeader(500)
	}
}
