package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"
	"github.com/dgmann/document-manager/pdf-processor/pdfprocessor"
	"net/http"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io/ioutil"
	"github.com/dgmann/document-manager/pdf-processor/image"
	"strconv"
)

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("images/convert", func(c *gin.Context) {
		images, err := pdfprocessor.ToImages(c.Request.Body)
		defer c.Request.Body.Close()
		if err != nil {
			c.Status(400)
			c.Error(err)
			log.Error(err)
			return
		}
		c.JSON(200, images)
	})
	router.POST("images/rotate/:degree", func(c *gin.Context) {
		degree, err := strconv.ParseFloat(c.Param("degree"), 64)
		if err != nil {
			c.AbortWithError(400, err)
		}

		img, err := ioutil.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()
		if err != nil || len(img) == 0 {
			c.AbortWithError(400, err)
		}
		rotated, format, err := image.Rotate(img, degree)
		if err != nil {
			c.AbortWithError(400, err)
		}
		c.Data(200, "image/"+format, rotated)
	})
	router.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "PDFProcessor")
	})
	router.Run(":8181")
}
