package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func registerRecords(g *gin.RouterGroup) {
	g.GET("", func(c *gin.Context) {
		r := c.Request.URL.Query()
		if _, ok := r["inbox"]; ok {
			records := app.Records.GetInbox()
			RespondAsJSON(c, records)
		} else {
			query := make(map[string]interface{})
			for k, v := range r {
				query[k] = v[0]
			}
			records := app.Records.Query(query)
			RespondAsJSON(c, records)
		}
	})

	g.GET("/:recordId", func(c *gin.Context) {
		id := c.Param("recordId")
		if id == "inbox" {
			records := app.Records.GetInbox()
			RespondAsJSON(c, records)
		} else {
			record := app.Records.Find(id)
			RespondAsJSON(c, record)
		}
	})

	g.GET("/:recordId/images/:imageId", func(c *gin.Context) {
		app.Images.Serve(c, c.Param("recordId"), c.Param("imageId"))
	})

	g.POST("", func(c *gin.Context) {
		file, _ := c.FormFile("pdf")
		sender := c.PostForm("sender")
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
		images, err := app.PDFProcessor.Convert(f)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		record, err := app.Records.Create(sender, images)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		c.Status(201)
		RespondAsJSON(c, record)
	})

	g.DELETE("/:recordId", func(c *gin.Context) {
		err := app.Records.Delete(c.Param("recordId"))
		c.Header("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			c.AbortWithError(400, err)
		}
		c.Status(204)
	})

	g.PATCH("/:recordId", func(c *gin.Context) {
		var record models.Record

		if err := json.NewDecoder(c.Request.Body).Decode(&record); err != nil {
			c.Error(err)
		}
		r := app.Records.Update(c.Param("recordId"), record)
		RespondAsJSON(c, r)
	})

	g.POST("/:recordId/append/:idtoappend", func(c *gin.Context) {
		recordToAppend := app.Records.Find(c.Param("idtoappend"))
		record := app.Records.Find(c.Param("recordId"))
		pages := append(record.Pages, recordToAppend.Pages...)

		err := app.Images.Copy(c.Param("idtoappend"), c.Param("recordId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		r := app.Records.Update(c.Param("recordId"), models.Record{Pages: pages})
		RespondAsJSON(c, r)
	})
}