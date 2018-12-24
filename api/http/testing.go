package http

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"net/url"
)

func NewTestContext() (*gin.Context, *gin.TestResponseRecorder) {
	w := gin.CreateTestResponseRecorder()
	context, _ := gin.CreateTestContext(w)
	context.Request = httptest.NewRequest("POST", "/", nil)
	requestUrl, _ := url.Parse("http://localhost")
	context.Set("location", requestUrl)
	return context, w
}
