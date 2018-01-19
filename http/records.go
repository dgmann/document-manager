package http

import (
	"github.com/dgmann/document-manager-api/services"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	log "github.com/sirupsen/logrus"
	"path"
	"strconv"
)

func registerRecords(g *gin.RouterGroup, recordDir string) {
	g.GET("", func(c *gin.Context) {
		r := records.GetInbox()
		if err := jsonapi.MarshalPayload(c.Writer, r); err != nil {
			c.Error(err)
		}
	})

	g.GET("/:recordId", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("recordId"), 10, 64)
		if err != nil {
			c.Error(jsonapi.ErrBadJSONAPIID)
		}
		record := records.Find(id)
		if err := jsonapi.MarshalPayload(c.Writer, record); err != nil {
			c.Error(err)
		}
	})

	g.GET("/:recordId/images/:imageId", func(c *gin.Context) {
		p := path.Join(recordDir, c.Param("recordId"), c.Param("imageId")+".png")
		c.File(p)
	})

	g.POST("", func(c *gin.Context) {
		file, _ := c.FormFile("pdf")
		f, err := file.Open()
		defer f.Close()
		if err != nil {
			fields := log.Fields{
				"name":  file.Filename,
				"size":  file.Size,
				"error": err,
			}
			log.WithFields(fields).Panic("Error opening PDF")
		}

		pdfProcessor := services.NewPDFProcessor("http://10.0.0.38:8181")
		images := pdfProcessor.ToImages(f)
		log.Debugf("Fetched %d images", len(images))

		sender := c.PostForm("sender")
		record := records.Create(sender)
		c.Status(201)
		c.Header("Content-Type", "application/json; charset=utf-8")
		if err := jsonapi.MarshalPayload(c.Writer, record); err != nil {
			c.Error(err)
		}
	})
}
