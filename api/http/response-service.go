package http

import (
	"github.com/dgmann/document-manager/api/models"
	"github.com/dgmann/document-manager/api/services"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseService struct {
	fileInfoService services.FileInfoService
}

func NewResponseService(fileInfoService services.FileInfoService) *ResponseService {
	return &ResponseService{fileInfoService: fileInfoService}
}

func (r *ResponseService) NewResponse(c *gin.Context, data interface{}) *Response {
	return r.NewResponseWithStatus(c, data, http.StatusOK)
}

func (r *ResponseService) NewResponseWithStatus(c *gin.Context, data interface{}, code int) *Response {
	baseUrl := ""
	if c != nil {
		url := location.Get(c)
		baseUrl = url.String()
	}
	switch data.(type) {
	case *models.Record:
		services.SetURL(data.(*models.Record), baseUrl, r.fileInfoService)
	case []*models.Record:
		for _, m := range data.([]*models.Record) {
			services.SetURL(m, baseUrl, r.fileInfoService)
		}
	}
	return &Response{Data: data, context: c, StatusCode: code}
}

func (r *ResponseService) NewErrorResponse(c *gin.Context, code int, err error) *Response {
	return &Response{Data: gin.H{"error": err.Error()}, context: c, StatusCode: code}
}
