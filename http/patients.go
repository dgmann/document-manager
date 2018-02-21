package http

import (
	"github.com/gin-gonic/gin"
)

func registerPatients(g *gin.RouterGroup) {
	g.GET("/:patientId/tags", func(c *gin.Context) {
		tags, err := app.Tags.ByPatient(c.Param("patientId"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		RespondAsJSON(c, tags)
	})
}
