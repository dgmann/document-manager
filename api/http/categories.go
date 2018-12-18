package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/models"
	"github.com/dgmann/document-manager/api/repositories/category"
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerCategories(g *gin.RouterGroup, factory Factory) {
	categoryController := NewCategoryController(factory)

	g.GET("", categoryController.All)
	g.POST("", categoryController.Create)
}

type CategoryController struct {
	categories      category.Repository
	responseService *ResponseService
}

func NewCategoryController(factory Factory) *CategoryController {
	return &CategoryController{
		categories:      factory.GetCategoryRepository(),
		responseService: factory.GetResponseService(),
	}
}

func (cat *CategoryController) All(c *gin.Context) {
	categories, err := cat.categories.All()
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	response := cat.responseService.NewResponse(c, categories)
	response.JSON()
}

func (cat *CategoryController) Create(c *gin.Context) {
	var body models.Category
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		cat.responseService.NewErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := cat.categories.Add(body.Id, body.Name); err != nil {
		cat.responseService.NewErrorResponse(c, http.StatusConflict, err)
		return
	}
	response := cat.responseService.NewResponse(c, body)
	response.JSON()
}
