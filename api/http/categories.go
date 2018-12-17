package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerCategories(g *gin.RouterGroup, factory Factory) {
	categoryRepository := factory.GetCategoryRepository()
	responseService := factory.GetResponseService()

	g.GET("", func(c *gin.Context) {
		cat, err := categoryRepository.All()
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		response := responseService.NewResponse(c, cat)
		response.JSON()
	})

	g.POST("", func(c *gin.Context) {
		var category models.Category
		if err := json.NewDecoder(c.Request.Body).Decode(&category); err != nil {
			responseService.NewErrorResponse(c, http.StatusBadRequest, err)
			return
		}
		if err := categoryRepository.Add(category.Id, category.Name); err != nil {
			responseService.NewErrorResponse(c, http.StatusConflict, err)
			return
		}
		response := responseService.NewResponse(c, category)
		response.JSON()
	})
}
