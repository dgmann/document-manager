package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/app"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseFactory struct {
	reader app.ModTimeReader
}

func NewResponseFactory(reader app.ModTimeReader) *ResponseFactory {
	return &ResponseFactory{reader: reader}
}

func (r *ResponseFactory) NewResponse(c *gin.Context, data interface{}) *Response {
	return r.NewResponseWithStatus(c, data, http.StatusOK)
}

func (r *ResponseFactory) NewResponseWithStatus(c *gin.Context, data interface{}, code int) *Response {
	baseUrl := ""
	if c != nil {
		url := location.Get(c)
		baseUrl = url.String()
	}
	payload := SetURL(data, baseUrl, r.reader)
	return &Response{Data: payload, context: c, StatusCode: code}
}

func (r *ResponseFactory) NewErrorResponse(c *gin.Context, code int, err error) *Response {
	return &Response{Data: gin.H{"error": err.Error()}, context: c, StatusCode: code}
}

type Response struct {
	StatusCode int
	Data       interface{}
	context    *gin.Context
}

func (r *Response) JSON() {
	r.context.Status(r.StatusCode)
	r.context.Header("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(r.context.Writer).Encode(r.Data); err != nil {
		r.context.AbortWithError(400, err)
	}
}
