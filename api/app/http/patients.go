package http

import (
	"github.com/dgmann/document-manager/api/app"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
)

func registerPatients(g *gin.RouterGroup, controller *PatientController) {
	g.GET("/:patientId/tags", controller.Tags)
	g.GET("/:patientId/categories", controller.Categories)
	g.GET("/:patientId/records", controller.Records)
}

type PatientController struct {
	records         app.RecordService
	tags            app.TagService
	categories      app.CategoryService
	responseService Responder
}

func (p *PatientController) Tags(c *gin.Context) {
	patientId := c.Param("patientId")

	tags, err := p.tags.ByPatient(c, patientId)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	res := p.responseService.NewResponse(c, tags)
	res.JSON()
}

func (p *PatientController) Categories(c *gin.Context) {
	patientId := c.Param("patientId")

	categories, err := p.categories.FindByPatient(c, patientId)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	res := p.responseService.NewResponse(c, categories)
	res.JSON()
}

func (p *PatientController) Records(c *gin.Context) {
	id := c.Param("patientId")
	records, err := p.records.Query(c, bson.M{"patientId": id})
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	res := p.responseService.NewResponse(c, records)
	res.JSON()
}
