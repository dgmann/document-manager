package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strconv"
	"errors"
	"bytes"
	"github.com/globalsign/mgo/bson"
)

func registerRecords(g *gin.RouterGroup) {
	g.GET("", func(c *gin.Context) {
		r := c.Request.URL.Query()
		var records []*models.Record
		var err error
		if _, ok := r["inbox"]; ok {
			records, err = app.Records.Query(bson.M{"$or": []bson.M{{"date": nil}, {"patientId": ""}, {"categoryId": nil}}})
		} else {
			query := make(map[string]interface{})
			for k, v := range r {
				query[k] = v[0]
			}
			records, err = app.Records.Query(query)
		}
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		RespondAsJSON(c, records)
	})

	g.GET("/:recordId", func(c *gin.Context) {
		id := c.Param("recordId")
		record, err := app.Records.Find(id)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		RespondAsJSON(c, record)
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
			return
		}
		c.Status(204)
	})

	g.PATCH("/:recordId", func(c *gin.Context) {
		var record models.Record

		if err := json.NewDecoder(c.Request.Body).Decode(&record); err != nil {
			c.Error(err)
			return
		}
		r, err := app.Records.Update(c.Param("recordId"), record)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		RespondAsJSON(c, r)
	})

	g.POST("/:recordId/append/:idtoappend", func(c *gin.Context) {
		recordToAppend, err := app.Records.Find(c.Param("idtoappend"))
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		record, err := app.Records.Find(c.Param("recordId"))
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		pages := append(record.Pages, recordToAppend.Pages...)

		err = app.Images.Copy(c.Param("idtoappend"), c.Param("recordId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		r, err := app.Records.Update(c.Param("recordId"), models.Record{Pages: pages})
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		RespondAsJSON(c, r)
	})

	g.GET("/:recordId/pages/:imageId", func(c *gin.Context) {
		record, err := app.Records.Find(c.Param("recordId"))
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		for _, page := range record.Pages {
			if page.Id == c.Param("imageId") {
				app.Images.Serve(c, c.Param("recordId"), c.Param("imageId"), page.Format)
				return
			}
		}
		c.AbortWithError(404, errors.New("page not found"))
	})

	g.POST("/:recordId/pages", func(c *gin.Context) {
		var updates []*models.PageUpdate
		if err := json.NewDecoder(c.Request.Body).Decode(&updates); err != nil {
			c.Error(err)
			return
		}

		r, err := app.Records.UpdatePages(c.Param("recordId"), updates)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		RespondAsJSON(c, r)
	})

	g.POST("/:recordId/pages/:imageId/rotate/:degrees", func(c *gin.Context) {
		images, err := app.Images.Get(c.Param("recordId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		degrees, err := strconv.Atoi(c.Param("degrees"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		if img, ok := images[c.Param("imageId")]; ok {
			img, err := app.PDFProcessor.Rotate(bytes.NewBuffer(img.Image), degrees)
			if err != nil {
				c.AbortWithError(400, err)
				return
			}
			app.Images.SetImage(c.Param("recordId"), c.Param("imageId"), img)
			c.JSON(200, img)
		} else {
			c.AbortWithError(400, errors.New("cannot read image"))
			return
		}
	})
}
