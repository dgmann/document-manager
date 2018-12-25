package http

import (
	"github.com/dgmann/document-manager/api/http/response"
	"github.com/dgmann/document-manager/api/repositories/category"
	"github.com/dgmann/document-manager/api/repositories/record"
	"github.com/dgmann/document-manager/api/repositories/tag"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
)

func registerPatients(g *gin.RouterGroup, factory Factory) {
	patientController := NewPatientController(factory)

	g.GET("/:patientId/tags", patientController.Tags)
	g.GET("/:patientId/categories", patientController.Categories)
	g.GET("/:patientId/records", patientController.Records)
}

type PatientController struct {
	records         record.Repository
	tags            tag.Repository
	categories      category.Repository
	responseService *response.Factory
}

func NewPatientController(factory Factory) *PatientController {
	return &PatientController{
		records:         factory.GetRecordRepository(),
		tags:            factory.GetTagRepository(),
		categories:      factory.GetCategoryRepository(),
		responseService: factory.GetResponseService(),
	}
}

func (p *PatientController) Tags(c *gin.Context) {
	patientId := c.Param("patientId")

	tags, err := p.tags.ByPatient(c.Request.Context(), patientId)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	res := p.responseService.NewResponse(c, tags)
	res.JSON()
}

func (p *PatientController) Categories(c *gin.Context) {
	patientId := c.Param("patientId")

	categories, err := p.categories.FindByPatient(c.Request.Context(), patientId)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	res := p.responseService.NewResponse(c, categories)
	res.JSON()
}

func (p *PatientController) Records(c *gin.Context) {
	id := c.Param("patientId")
	records, err := p.records.Query(c.Request.Context(), bson.M{"patientId": id})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	res := p.responseService.NewResponse(c, records)
	res.JSON()
}
