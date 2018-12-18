package http

import (
	"github.com/dgmann/document-manager/api/http/response"
	"github.com/gin-gonic/gin"
)

type Responder interface {
	OkResponder
	ErrorResponder
	StatusResponder
}

type OkResponder interface {
	NewResponse(c *gin.Context, data interface{}) *response.Response
}

type ErrorResponder interface {
	NewErrorResponse(c *gin.Context, code int, err error) *response.Response
}

type StatusResponder interface {
	NewResponseWithStatus(c *gin.Context, data interface{}, code int) *response.Response
}
