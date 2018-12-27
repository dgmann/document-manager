package http

import (
	"github.com/gin-gonic/gin"
)

type Responder interface {
	OkResponder
	ErrorResponder
	StatusResponder
}

type OkResponder interface {
	NewResponse(c *gin.Context, data interface{}) *Response
}

type ErrorResponder interface {
	NewErrorResponse(c *gin.Context, code int, err error) *Response
}

type StatusResponder interface {
	NewResponseWithStatus(c *gin.Context, data interface{}, code int) *Response
}
