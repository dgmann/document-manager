package http

import (
	"github.com/gin-gonic/gin"
	"net/url"
)

func NewTestContext() (*gin.Context, *gin.TestResponseRecorder) {
	w := gin.CreateTestResponseRecorder()
	context, _ := gin.CreateTestContext(w)
	requestUrl, _ := url.Parse("http://localhost")
	context.Set("location", requestUrl)
	return context, w
}
