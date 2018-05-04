package http

import (
	"github.com/gin-gonic/gin"
	"encoding/json"
	"github.com/dgmann/document-manager/api/models"
)

func registerCategories(g *gin.RouterGroup, factory *Factory) {
	categoryRepository := factory.GetCategoryRepository()

	g.GET("", func(c *gin.Context) {
		cat, err := categoryRepository.All()
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		RespondAsJSON(c, cat)
	})

	g.POST("", func(c *gin.Context) {
		var category models.Category
		if err := json.NewDecoder(c.Request.Body).Decode(&category); err != nil {
			c.Error(err)
			c.AbortWithError(400, err)
			return
		}
		if err := categoryRepository.Add(category.Name); err != nil {
			c.Error(err)
			c.AbortWithError(400, err)
			return
		}
		RespondAsJSON(c, category)
	})
}
