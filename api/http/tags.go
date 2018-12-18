package http

import (
	"github.com/dgmann/document-manager/api/repositories/tag"
	"github.com/gin-gonic/gin"
)

type TagController struct {
	tags tag.Repository
}

type tagControllerFactory interface {
	GetTagRepository() tag.Repository
}

func NewTagController(factory tagControllerFactory) *TagController {
	return &TagController{tags: factory.GetTagRepository()}
}

func (t *TagController) All(c *gin.Context) {
	tags, err := t.tags.All()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, tags)
}
