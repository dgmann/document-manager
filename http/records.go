package http

import (
	"github.com/dgmann/document-manager-api/models"
	"github.com/dgmann/document-manager-api/services"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	log "github.com/sirupsen/logrus"
	"path"
	"time"
)

func registerRecords(g *gin.RouterGroup, recordDir string) {
	g.GET("", func(c *gin.Context) {
		records := []*models.Record{getData(c.Request.Host)}
		if err := jsonapi.MarshalPayload(c.Writer, records); err != nil {
			c.Error(err)
		}
	})

	g.GET("/:recordId", func(c *gin.Context) {
		record := getData(c.Request.Host)
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

		pdfProcessor := services.NewPDFProcessor("http://localhost:8181")
		images := pdfProcessor.ToImages(f)
		c.JSON(200, images)
	})
}

func getData(url string) *models.Record {
	pages := []models.Page{
		{Index: 0, Content: "", Url: "http://" + path.Join(url, "/records/1/images/1234")},
		{Index: 1, Content: "", Url: "http://" + path.Join(url, "/records/1/images/1234")},
		{Index: 2, Content: "", Url: "http://" + path.Join(url, "/records/1/images/quer")},
		{Index: 3, Content: "", Url: "http://" + path.Join(url, "/records/1/images/1234")},
	}
	return &models.Record{Id: "1", Date: time.Now(), Comment: "Neuer Patient? Was sollen wir mit dem Bazi machen?", Sender: "Scan", Pages: pages}
}
