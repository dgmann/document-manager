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
		file, err := c.FormFile("pdf")
		if err != nil {
			fields := log.Fields{
				"name": file.Filename,
				"size": file.Size,
				"error": err,
			}
			log.WithFields(fields).Panic("Error reading form field")
		}

		f, err := file.Open()
		defer f.Close()
		if err != nil {
			fields := log.Fields{
				"name": file.Filename,
				"size": file.Size,
				"error": err,
			}
			log.WithFields(fields).Panic("Error opening PDF")
		}


		images := pdfprocessor.ToImages(f)
		c.JSON(200, images)
	})
	router.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "PDFProcessor")
	})
	router.Run(":8181")
}
