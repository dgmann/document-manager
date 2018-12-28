package http

import (
	"context"
	"encoding/json"
	"github.com/dgmann/document-manager/api/app"
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerCategories(g *gin.RouterGroup, controller *CategoryController) {
	g.GET("", controller.All)
	g.POST("", controller.Create)
}

type CategoryController struct {
	categories      categoryRepository
	responseService Responder
}

type categoryRepository interface {
	All(ctx context.Context) ([]app.Category, error)
	Add(ctx context.Context, id, category string) error
}

func (cat *CategoryController) All(c *gin.Context) {
	categories, err := cat.categories.All(c)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	resp := cat.responseService.NewResponse(c, categories)
	resp.JSON()
}

func (cat *CategoryController) Create(c *gin.Context) {
	var body app.Category
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		cat.responseService.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := cat.categories.Add(c.Request.Context(), body.Id, body.Name); err != nil {
		cat.responseService.NewErrorResponse(c, http.StatusConflict, err)
		return
	}
	resp := cat.responseService.NewResponseWithStatus(c, body, 201)
	resp.JSON()
}
