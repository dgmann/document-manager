package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/http/response"
	"github.com/dgmann/document-manager/api/models"
	"github.com/dgmann/document-manager/api/repositories/category"
	"github.com/dgmann/document-manager/api/repositories/patient"
	"github.com/dgmann/document-manager/api/repositories/record"
	"github.com/dgmann/document-manager/api/repositories/tag"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"strings"
)

func registerPatients(g *gin.RouterGroup, factory Factory) {
	patientController := NewPatientController(factory)

	g.GET("", patientController.All)
	g.POST("", patientController.Create)
	g.GET("/:patientId", patientController.One)

	g.GET("/:patientId/tags", patientController.Tags)
	g.GET("/:patientId/categories", patientController.Categories)
	g.GET("/:patientId/records", patientController.Records)
}

type PatientController struct {
	records         record.Repository
	tags            tag.Repository
	patients        patient.Repository
	categories      category.Repository
	responseService *response.Factory
}

func NewPatientController(factory Factory) *PatientController {
	return &PatientController{
		records:         factory.GetRecordRepository(),
		tags:            factory.GetTagRepository(),
		patients:        factory.GetPatientRepository(),
		categories:      factory.GetCategoryRepository(),
		responseService: factory.GetResponseService(),
	}
}

func (p *PatientController) All(c *gin.Context) {
	var patients []*models.Patient
	var err error
	name := c.DefaultQuery("name", "")
	if len(name) > 0 {
		names := strings.Split(name, ",")
		lastName, firstName := names[0], ""
		if len(names) > 1 {
			firstName = names[1]
		}
		patients, err = p.patients.FindByName(firstName+".*", lastName+".*")
	} else {
		patients, err = p.patients.All()
	}

	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	response := p.responseService.NewResponse(c, patients)
	response.JSON()
}

func (p *PatientController) One(c *gin.Context) {
	patientId := c.Param("patientId")

	tags, err := p.tags.ByPatient(patientId)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	categories, err := p.categories.FindByPatient(patientId)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	result, err := p.patients.Find(patientId)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	result.Tags = tags
	result.Categories = categories
	response := p.responseService.NewResponse(c, result)
	response.JSON()
}

func (p *PatientController) Create(c *gin.Context) {
	var body models.Patient
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		c.Error(err)
		c.AbortWithError(400, err)
		return
	}
	if err := p.patients.Add(&body); err != nil {
		c.Error(err)
		c.AbortWithError(400, err)
		return
	}
	response := p.responseService.NewResponse(c, body)
	response.JSON()
}

func (p *PatientController) Tags(c *gin.Context) {
	patientId := c.Param("patientId")

	tags, err := p.tags.ByPatient(patientId)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	response := p.responseService.NewResponse(c, tags)
	response.JSON()
}

func (p *PatientController) Categories(c *gin.Context) {
	patientId := c.Param("patientId")

	categories, err := p.categories.FindByPatient(patientId)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	response := p.responseService.NewResponse(c, categories)
	response.JSON()
}

func (p *PatientController) Records(c *gin.Context) {
	id := c.Param("patientId")
	records, err := p.records.Query(bson.M{"patientId": id})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	response := p.responseService.NewResponse(c, records)
	response.JSON()
}
