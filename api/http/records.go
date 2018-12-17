package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/api/models"
	"github.com/dgmann/document-manager/api/repositories"
	"github.com/dgmann/document-manager/api/repositories/image"
	"github.com/dgmann/document-manager/api/repositories/pdf"
	"github.com/dgmann/document-manager/api/repositories/record"
	"github.com/dgmann/document-manager/api/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

func registerRecords(g *gin.RouterGroup, factory Factory) {
	records := NewRecordController(factory)

	g.GET("", records.All)
	g.GET("/:recordId", records.One)
	g.POST("", records.Create)
	g.PATCH("/:recordId", records.Update)
	g.DELETE("/:recordId", records.Delete)

	g.POST("/:recordId/duplicate", records.Duplicate)
	g.PUT("/:recordId/reset", records.Reset)
	g.POST("/:recordId/append/:idtoappend", records.Append)

	g.GET("/:recordId/pages/:imageId", records.Page)
	g.POST("/:recordId/pages", records.UpdatePages)
	g.POST("/:recordId/pages/:imageId/rotate/:degrees", records.RotatePage)
}

type RecordController struct {
	records         record.Repository
	images          image.Repository
	pdfs            pdf.Repository
	pdfProcessor    *services.PdfProcessor
	responseService *ResponseService
}

func NewRecordController(factory Factory) *RecordController {
	recordRepository := factory.GetRecordRepository()
	imageRepository := factory.GetImageRepository()
	pdfRepository := factory.GetPDFRepository()
	pdfProcessor, err := factory.GetPdfProcessor()
	if err != nil {
		log.WithError(err).Error("Cannot reach PDF processor")
	}
	responseService := factory.GetResponseService()

	return &RecordController{
		records:         recordRepository,
		images:          imageRepository,
		pdfs:            pdfRepository,
		pdfProcessor:    pdfProcessor,
		responseService: responseService,
	}
}

func (r *RecordController) All(c *gin.Context) {
	params := c.Request.URL.Query()
	query := make(map[string]interface{})
	for k, v := range params {
		query[k] = v[0]
	}
	records, err := r.records.Query(query)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	response := r.responseService.NewResponse(c, records)
	response.JSON()
}

func (r *RecordController) One(c *gin.Context) {
	id := c.Param("recordId")
	result, err := r.records.Find(id)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	response := r.responseService.NewResponse(c, result)
	response.JSON()
}

func (r *RecordController) Create(c *gin.Context) {
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

	images, err := r.pdfProcessor.Convert(bytes.NewBuffer(fileBytes))
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	res, err := r.records.Create(newRecord, images, bytes.NewBuffer(fileBytes))
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	response := r.responseService.NewResponseWithStatus(c, res, 201)
	response.JSON()
}

func (r *RecordController) Update(c *gin.Context) {
	var body models.Record

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		c.Error(err)
		return
	}
	updated, err := r.records.Update(c.Param("recordId"), body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	response := r.responseService.NewResponse(c, updated)
	response.JSON()
}

func (r *RecordController) Delete(c *gin.Context) {
	err := r.records.Delete(c.Param("recordId"))
	c.Header("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	c.Status(204)
}

func (r *RecordController) Duplicate(c *gin.Context) {
	recordToDuplicate, err := r.records.Find(c.Param("recordId"))
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	file, err := r.pdfs.Get(c.Param("recordId"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	newId := bson.NewObjectId()

	err = r.images.Copy(recordToDuplicate.Id.Hex(), newId.Hex())
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	copiedRecord, err := r.records.Create(models.CreateRecord{
		Id:         &newId,
		ReceivedAt: recordToDuplicate.ReceivedAt,
		Sender:     recordToDuplicate.Sender,
		Pages:      recordToDuplicate.Pages,
	}, nil, file)

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	response := r.responseService.NewResponse(c, copiedRecord)
	response.JSON()
}

func (r *RecordController) Reset(c *gin.Context) {
	recordId := c.Param("recordId")
	f, err := r.pdfs.Get(recordId)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	images, err := r.pdfProcessor.Convert(bytes.NewBuffer(fileBytes))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	var pages []*models.Page
	for _, img := range images {
		page := models.NewPage(img.Format)
		if err := r.images.Write(repositories.NewKeyedGenericResource(img.Image, img.Format, recordId, page.Id)); err != nil {
			c.AbortWithError(500, err)
			return
		}
		pages = append(pages, page)
	}

	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	updated, err := r.records.Update(recordId, models.Record{Pages: pages})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	response := r.responseService.NewResponse(c, updated)
	response.JSON()
}

func (r *RecordController) Append(c *gin.Context) {
	recordToAppend, err := r.records.Find(c.Param("idtoappend"))
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	targetRecord, err := r.records.Find(c.Param("recordId"))
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	pages := append(targetRecord.Pages, recordToAppend.Pages...)

	err = r.images.Copy(c.Param("idtoappend"), c.Param("recordId"))
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	updated, err := r.records.Update(c.Param("recordId"), models.Record{Pages: pages})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	response := r.responseService.NewResponse(c, updated)
	response.JSON()
}

func (r *RecordController) Page(c *gin.Context) {
	rec, err := r.records.Find(c.Param("recordId"))
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	for _, page := range rec.Pages {
		if page.Id == c.Param("imageId") {
			r.images.Serve(c, c.Param("recordId"), c.Param("imageId"), page.Format)
			return
		}
	}
	c.AbortWithError(404, errors.New("page not found"))
}

func (r *RecordController) UpdatePages(c *gin.Context) {
	var updates []*models.PageUpdate
	if err := json.NewDecoder(c.Request.Body).Decode(&updates); err != nil {
		c.Error(err)
		return
	}

	updated, err := r.records.UpdatePages(c.Param("recordId"), updates)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	images, err := r.images.Get(c.Param("recordId"))
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
				img, err := r.pdfProcessor.Rotate(bytes.NewBuffer(img.Image), int(update.Rotate))
				if err != nil {
					mutex.Lock()
					errorIds = append(errorIds, update.Id)
					mutex.Unlock()
					return
				}
				err = r.images.Write(repositories.NewKeyedGenericResource(img.Image, img.Format, c.Param("recordId"), update.Id))
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

	response := r.responseService.NewResponse(c, updated)
	response.JSON()
}

func (r *RecordController) RotatePage(c *gin.Context) {
	images, err := r.images.Get(c.Param("recordId"))
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
		img, err := r.pdfProcessor.Rotate(bytes.NewBuffer(img.Image), degrees)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		if err := r.images.Write(repositories.NewKeyedGenericResource(img.Image, img.Format, c.Param("recordId"), c.Param("imageId"))); err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, img)
	} else {
		c.AbortWithError(400, errors.New("cannot read image"))
		return
	}
}
