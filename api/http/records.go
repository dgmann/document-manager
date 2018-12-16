package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/models"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

func registerRecords(g *gin.RouterGroup, factory *Factory) {
	recordRepository := factory.GetRecordRepository()
	imageRepository := factory.GetImageRepository()
	pdfProcessor, err := factory.GetPdfProcessor()
	if err != nil {
		log.WithError(err).Error("Cannot reach PDF processor")
	}
	responseService := factory.GetResponseService()

	g.GET("", func(c *gin.Context) {
		r := c.Request.URL.Query()
		query := make(map[string]interface{})
		for k, v := range r {
			query[k] = v[0]
		}
		records, err := recordRepository.Query(query)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		response := responseService.NewResponse(c, records)
		response.JSON()
	})

	g.GET("/:recordId", func(c *gin.Context) {
		id := c.Param("recordId")
		record, err := recordRepository.Find(id)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		response := responseService.NewResponse(c, record)
		response.JSON()
	})

	g.POST("", func(c *gin.Context) {
		var newRecord models.CreateRecord
		if err := c.Bind(&newRecord); err != nil {
			log.WithError(err).Error("error decoding Data")
			c.AbortWithError(400, err)
			return
		}

		file, err := c.FormFile("pdf")
		if err != nil {
			log.WithError(err).Error("no file specified in upload")
			c.AbortWithError(400, err)
			return
		}

		f, err := file.Open()
		if err != nil {
			fields := log.Fields{
				"name":  file.Filename,
				"size":  file.Size,
				"error": err,
			}
			log.WithFields(fields).Panic("Error opening PDF")
			c.AbortWithError(400, err)
			return
		}
		defer f.Close()
		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		images, err := pdfProcessor.Convert(bytes.NewBuffer(fileBytes))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		record, err := recordRepository.Create(newRecord, images, bytes.NewBuffer(fileBytes))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		response := responseService.NewResponseWithStatus(c, record, 201)
		response.JSON()
	})

	g.DELETE("/:recordId", func(c *gin.Context) {
		err := recordRepository.Delete(c.Param("recordId"))
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
		r, err := recordRepository.Update(c.Param("recordId"), record)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		if record.PatientId != nil {

		}

		response := responseService.NewResponse(c, r)
		response.JSON()
	})

	g.POST("/:recordId/duplicate", func(c *gin.Context) {
		recordToDuplicate, err := recordRepository.Find(c.Param("recordId"))
		if err != nil {
			c.AbortWithError(404, err)
			return
		}
		pdfs := factory.GetPDFRepository()
		file, err := pdfs.Get(c.Param("recordId"))
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		newId := bson.NewObjectId()

		err = imageRepository.Copy(recordToDuplicate.Id.Hex(), newId.Hex())
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		copiedRecord, err := recordRepository.Create(models.CreateRecord{
			Id:         &newId,
			ReceivedAt: recordToDuplicate.ReceivedAt,
			Sender:     recordToDuplicate.Sender,
			Pages:      recordToDuplicate.Pages,
		}, nil, file)

		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		response := responseService.NewResponse(c, copiedRecord)
		response.JSON()
	})

	g.PUT("/:recordId/reset", func(c *gin.Context) {
		recordId := c.Param("recordId")
		pdfs := factory.GetPDFRepository()
		f, err := pdfs.Get(recordId)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		images, err := pdfProcessor.Convert(bytes.NewBuffer(fileBytes))
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		pages, err := imageRepository.Set(recordId, images)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		updated, err := recordRepository.Update(recordId, models.Record{Pages: pages})
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		response := responseService.NewResponse(c, updated)
		response.JSON()
	})

	g.POST("/:recordId/append/:idtoappend", func(c *gin.Context) {
		recordToAppend, err := recordRepository.Find(c.Param("idtoappend"))
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		record, err := recordRepository.Find(c.Param("recordId"))
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		pages := append(record.Pages, recordToAppend.Pages...)

		err = imageRepository.Copy(c.Param("idtoappend"), c.Param("recordId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		r, err := recordRepository.Update(c.Param("recordId"), models.Record{Pages: pages})
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		response := responseService.NewResponse(c, r)
		response.JSON()
	})

	g.GET("/:recordId/pages/:imageId", func(c *gin.Context) {
		record, err := recordRepository.Find(c.Param("recordId"))
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		for _, page := range record.Pages {
			if page.Id == c.Param("imageId") {
				imageRepository.Serve(c, c.Param("recordId"), c.Param("imageId"), page.Format)
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

		r, err := recordRepository.UpdatePages(c.Param("recordId"), updates)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		images, err := imageRepository.Get(c.Param("recordId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		var errorIds []string
		var wg sync.WaitGroup
		var mutex = &sync.Mutex{}

		for _, u := range updates {
			if u.Rotate == 0 {
				continue
			}

			wg.Add(1)
			go func(update *models.PageUpdate) {
				defer wg.Done()

				if img, ok := images[update.Id]; ok {
					img, err := pdfProcessor.Rotate(bytes.NewBuffer(img.Image), int(update.Rotate))
					if err != nil {
						mutex.Lock()
						errorIds = append(errorIds, update.Id)
						mutex.Unlock()
						return
					}
					imageRepository.SetImage(c.Param("recordId"), update.Id, img)
					if err != nil {
						mutex.Lock()
						errorIds = append(errorIds, update.Id)
						mutex.Unlock()
						return
					}
				} else {
					mutex.Lock()
					errorIds = append(errorIds, update.Id)
					mutex.Unlock()
				}
			}(u)
		}

		wg.Wait()
		if len(errorIds) > 0 {
			c.AbortWithError(503, errors.New(fmt.Sprintf("error rotating pages %s", strings.Join(errorIds, ","))))
			return
		}

		response := responseService.NewResponse(c, r)
		response.JSON()
	})

	g.POST("/:recordId/pages/:imageId/rotate/:degrees", func(c *gin.Context) {
		images, err := imageRepository.Get(c.Param("recordId"))
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
			img, err := pdfProcessor.Rotate(bytes.NewBuffer(img.Image), degrees)
			if err != nil {
				c.AbortWithError(400, err)
				return
			}
			imageRepository.SetImage(c.Param("recordId"), c.Param("imageId"), img)
			c.JSON(200, img)
		} else {
			c.AbortWithError(400, errors.New("cannot read image"))
			return
		}
	})
}
