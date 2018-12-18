package http

import (
	"github.com/dgmann/document-manager/api/services"
	"github.com/gin-gonic/gin"
)

type GeneralController struct {
	health services.HealthChecker
}

func NewGeneralController() *GeneralController {
	return &GeneralController{
		health: services.GetHealthService(),
	}
}

func (g *GeneralController) Home(c *gin.Context) {
	c.String(200, "Document Manager API")
}

func (g *GeneralController) Status(c *gin.Context) {
	if err := g.health.Check(); err != nil {
		c.String(500, "Status: Error, Message: %s", err)
	} else {
		c.String(200, "Status: Ok")
	}
}
