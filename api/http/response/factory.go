package response

import (
	"github.com/dgmann/document-manager/api/services"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Factory struct {
	reader services.ModTimeReader
}

func NewFactory(reader services.ModTimeReader) *Factory {
	return &Factory{reader: reader}
}

func (r *Factory) NewResponse(c *gin.Context, data interface{}) *Response {
	return r.NewResponseWithStatus(c, data, http.StatusOK)
}

func (r *Factory) NewResponseWithStatus(c *gin.Context, data interface{}, code int) *Response {
	baseUrl := ""
	if c != nil {
		url := location.Get(c)
		baseUrl = url.String()
	}
	payload := services.SetURL(data, baseUrl, r.reader)
	return &Response{Data: payload, context: c, StatusCode: code}
}

func (r *Factory) NewErrorResponse(c *gin.Context, code int, err error) *Response {
	return &Response{Data: gin.H{"error": err.Error()}, context: c, StatusCode: code}
}
