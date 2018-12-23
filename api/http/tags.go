package http

import (
	"github.com/dgmann/document-manager/api/repositories/tag"
	"github.com/gin-gonic/gin"
)

type TagController struct {
	tags tag.Repository
}

func NewTagController(repository tag.Repository) *TagController {
	return &TagController{tags: repository}
}

func (t *TagController) All(c *gin.Context) {
	tags, err := t.tags.All(c.Request.Context())
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, tags)
}
