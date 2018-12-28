package http

import (
	"github.com/dgmann/document-manager/api/app"
	"github.com/gin-gonic/gin"
)

type TagController struct {
	tags app.TagService
}

func NewTagController(repository app.TagService) *TagController {
	return &TagController{tags: repository}
}

func (t *TagController) All(c *gin.Context) {
	tags, err := t.tags.All(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, tags)
}
