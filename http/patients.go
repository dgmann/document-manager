package http

import (
	"github.com/gin-gonic/gin"
	"github.com/dgmann/document-manager-api/models"
)

func registerPatients(g *gin.RouterGroup) {

	g.GET("", func(c *gin.Context) {
		patients := []models.Patient{
			{Id: "1", Name: "John Doe"},
			{Id: "2", Name: "Donald Trump"},
			{Id: "3", Name: "Barack Obama"},
		}
		RespondAsJSON(c, patients)
	})

	g.GET("/:patientId/tags", func(c *gin.Context) {
		tags, err := app.Tags.ByPatient(c.Param("patientId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		RespondAsJSON(c, tags)
	})
}
