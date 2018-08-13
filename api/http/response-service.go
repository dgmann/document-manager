package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func RespondAsJSON(c *gin.Context, data interface{}) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(c.Writer).Encode(data); err != nil {
		c.AbortWithError(400, err)
	}
}

func RespondWithError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"error": err.Error()})
	return
}


