package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/dgmann/document-manager-api/models"
	"time"
	"path"
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
		path := path.Join(recordDir, c.Param("recordId"), c.Param("imageId") + ".png")
		c.File(path)
	})
}

func getData(url string) *models.Record {
	pages := []models.Page{
		{ Index: 0, Content: "", Url: "http://" + path.Join(url, "/records/1/images/1234") },
		{ Index: 1, Content: "", Url: "http://" + path.Join(url, "/records/1/images/1234") },
		{ Index: 2, Content: "", Url: "http://" + path.Join(url, "/records/1/images/quer") },
		{ Index: 3, Content: "", Url: "http://" + path.Join(url, "/records/1/images/1234") },
	}
	return &models.Record{Id: "1", Date: time.Now(), Comment: "Neu?", Sender: "Scan", Pages: pages}
}
