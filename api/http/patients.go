package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/models"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"strings"
)

func registerPatients(g *gin.RouterGroup, factory *Factory) {
	patientRepository := factory.GetPatientRepository()
	tagRepository := factory.GetTagRepository()
	recordRepository := factory.GetRecordRepository()
	categoryRepository := factory.GetCategoryRepository()
	responseService := factory.GetResponseService()

	g.POST("", func(c *gin.Context) {
		var patient models.Patient
		if err := json.NewDecoder(c.Request.Body).Decode(&patient); err != nil {
			c.Error(err)
			c.AbortWithError(400, err)
			return
		}
		if err := patientRepository.Add(&patient); err != nil {
			c.Error(err)
			c.AbortWithError(400, err)
			return
		}
		response := responseService.NewResponse(patient)
		RespondAsJSON(c, response)
	})

	g.GET("", func(c *gin.Context) {
		var patients []*models.Patient
		var err error
		name := c.DefaultQuery("name", "")
		if len(name) > 0 {
			names := strings.Split(name, ",")
			lastName, firstName := names[0], ""
			if len(names) > 1 {
				firstName = names[1]
			}
			patients, err = patientRepository.FindByName(firstName+".*", lastName+".*")
		} else {
			patients, err = patientRepository.All()
		}

		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		response := responseService.NewResponse(patients)
		RespondAsJSON(c, response)
	})

	g.GET("/:patientId", func(c *gin.Context) {
		patientId := c.Param("patientId")

		tags, err := tagRepository.ByPatient(patientId)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		categories, err := categoryRepository.FindByPatient(patientId)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		patient, err := patientRepository.Find(patientId)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}
		patient.Tags = tags
		patient.Categories = categories
		response := responseService.NewResponse(patient)
		RespondAsJSON(c, response)
	})

	g.GET("/:patientId/tags", func(c *gin.Context) {
		patientId := c.Param("patientId")

		tags, err := tagRepository.ByPatient(patientId)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		response := responseService.NewResponse(tags)
		RespondAsJSON(c, response)
	})

	g.GET("/:patientId/categories", func(c *gin.Context) {
		patientId := c.Param("patientId")

		categories, err := categoryRepository.FindByPatient(patientId)
		if err != nil {
			c.AbortWithError(404, err)
			return
		}

		response := responseService.NewResponse(categories)
		RespondAsJSON(c, response)
	})

	g.GET("/:patientId/records", func(c *gin.Context) {
		id := c.Param("patientId")
		records, err := recordRepository.Query(bson.M{"patientId": id})
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		response := responseService.NewResponse(records)
		RespondAsJSON(c, response)
	})
}
