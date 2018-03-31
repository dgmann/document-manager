package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"
	"github.com/dgmann/pdf-processor/pdfprocessor"
	"net/http"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("images/extract", func(c *gin.Context) {
		images, err := pdfprocessor.ToImages(c.Request.Body)
		if err != nil {
			c.Status(400)
			c.Error(err)
			log.Error(err)
			return
		}
		c.JSON(200, images)
	})
	router.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "PDFProcessor")
	})
	router.Run(":8181")
}
