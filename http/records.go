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
		pages := []models.Page{
			{ Index: 0, Content: "", Url: "http://" + path.Join(c.Request.Host, "/records/1/images/1234") },
			{ Index: 1, Content: "", Url: "http://" + path.Join(c.Request.Host, "/records/1/images/1234") },
		}
		records := []*models.Record{{Id: "1", Date: time.Now(), Comment: "New?", Sender: "Scan", Pages: pages}}
		if err := jsonapi.MarshalPayload(c.Writer, records); err != nil {
			c.Error(err)
		}
	})

	g.GET("/:recordId/images/:imageId", func(c *gin.Context) {
		path := path.Join(recordDir, c.Param("recordId"), c.Param("imageId") + ".png")
		c.File(path)
	})
}
