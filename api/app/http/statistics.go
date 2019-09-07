package http

import (
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"github.com/gin-gonic/gin"
)

type providers map[string]app.StatisticProvider

type StatisticsController struct {
	providers providers
}

func (h *StatisticsController) Statistics(c *gin.Context) {
	messages := make(map[string]string)
	for key, provider := range h.providers {
		numberOfElements, err := provider.NumberOfElements()
		if err != nil {
			messages[key] = err.Error()
		} else {
			messages[key] = fmt.Sprintf("number of elements: %d", numberOfElements)
		}
	}
	c.JSON(200, messages)
}
