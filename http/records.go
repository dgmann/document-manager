package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager-api/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/google/jsonapi"
	log "github.com/sirupsen/logrus"
	"path"
)

func registerRecords(g *gin.RouterGroup, recordDir string) {
	g.GET("", func(c *gin.Context) {
		r := records.GetInbox()
		c.Header("Content-Type", "application/json; charset=utf-8")
		if err := jsonapi.MarshalPayload(c.Writer, r); err != nil {
			c.Error(err)
		}
	})

	g.GET("/:recordId", func(c *gin.Context) {
		id := c.Param("recordId")
		record := records.Find(bson.ObjectIdHex(id))
		c.Header("Content-Type", "application/json; charset=utf-8")
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

	g.PUT("/:recordId", func(c *gin.Context) {
		var data interface{}

		if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
			c.Error(err)
		}
		record := records.Update(c.Param("recordId"), data.(map[string]interface{})["data"].(map[string]interface{})["attributes"].(map[string]interface{}))
		c.Header("Content-Type", "application/json; charset=utf-8")
		if err := jsonapi.MarshalPayload(c.Writer, record); err != nil {
			c.Error(err)
		}
	})
}
