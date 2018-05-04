package http

import (
	"github.com/gin-gonic/gin"
	"github.com/dgmann/document-manager/api/models"
	"strings"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
)

func registerPatients(g *gin.RouterGroup, factory *Factory) {
	patientRepository := factory.GetPatientRepository()
	tagRepository := factory.GetTagRepository()
	recordRepository := factory.GetRecordRepository()
	categoryRepository := factory.GetCategoryRepository()

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
		RespondAsJSON(c, patient)
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
		RespondAsJSON(c, patients)
	})

	g.GET("/:patientId", func(c *gin.Context) {
		tags, err := tagRepository.ByPatient(c.Param("patientId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		categories, err := categoryRepository.FindByPatient(c.Param("patientId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		patient, err := patientRepository.Find(c.Param("patientId"))
		if err != nil {
			c.AbortWithError(404, err)
			return
		}
		patient.Tags = tags
		patient.Categories = categories
		RespondAsJSON(c, patient)
	})

	g.GET("/:patientId/records", func(c *gin.Context) {
		id := c.Param("patientId")
		records, err := recordRepository.Query(bson.M{"patientId": id})
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		RespondAsJSON(c, records)
	})
}
