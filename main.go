package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("images/extract", func(c *gin.Context) {
		file, _ := c.FormFile("pdf")
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


		images := toImages(f)
		c.JSON(200, images)
	})
	router.Run()
}
