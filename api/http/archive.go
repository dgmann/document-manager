package http

import (
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
)

type ArchiveController struct {
	pdfs getter
}

type getter interface {
	Get(id string) (io.Reader, error)
}

func NewArchiveController(pdf getter) *ArchiveController {
	return &ArchiveController{pdfs: pdf}
}

func (a *ArchiveController) One(c *gin.Context) {
	file, err := a.pdfs.Get(c.Param("recordId"))
	if err != nil {
		c.AbortWithError(404, err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Data(200, "application/pdf", data)
}
