package http

import (
	"context"
	"encoding/json"
	"github.com/dgmann/document-manager/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerCategories(g *gin.RouterGroup, factory Factory) {
	categoryController := NewCategoryController(factory.GetCategoryRepository(), factory.GetResponseService())

	g.GET("", categoryController.All)
	g.POST("", categoryController.Create)
}

type CategoryController struct {
	categories      categoryRepository
	responseService Responder
}

type categoryRepository interface {
	All(ctx context.Context) ([]models.Category, error)
	Add(ctx context.Context, id, category string) error
}

func NewCategoryController(categories categoryRepository, responseService Responder) *CategoryController {
	return &CategoryController{
		categories:      categories,
		responseService: responseService,
	}
}

func (cat *CategoryController) All(c *gin.Context) {
	categories, err := cat.categories.All(c.Request.Context())
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	resp := cat.responseService.NewResponse(c, categories)
	resp.JSON()
}

func (cat *CategoryController) Create(c *gin.Context) {
	var body models.Category
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
