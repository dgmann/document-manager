package http

import (
	"github.com/gin-gonic/gin"
	"github.com/dgmann/document-manager/api/models"
	"strings"
	"encoding/json"
)

func registerPatients(g *gin.RouterGroup) {

	g.POST("", func(c *gin.Context) {
		var patient models.Patient
		if err := json.NewDecoder(c.Request.Body).Decode(&patient); err != nil {
			c.Error(err)
			c.AbortWithError(400, err)
			return
		}
		if err := app.Patients.Add(&patient); err != nil {
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
			patients, err = app.Patients.FindByName(firstName+".*", lastName+".*")
		} else {
			patients, err = app.Patients.All()
		}

		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		RespondAsJSON(c, patients)
	})

	g.GET("/:patientId", func(c *gin.Context) {
		tags, err := app.Tags.ByPatient(c.Param("patientId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		categories, err := app.Categories.FindByPatient(c.Param("patientId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		patient, err := app.Patients.Find(c.Param("patientId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		patient.Tags = tags
		patient.Categories = categories
		RespondAsJSON(c, patient)
	})

	g.GET("/:patientId/records", func(c *gin.Context) {
		records, err := app.Records.FindByPatientId(c.Param("patientId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		RespondAsJSON(c, records)
	})
}
