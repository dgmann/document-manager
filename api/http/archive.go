package http

import (
	"github.com/dgmann/document-manager/api/repositories/pdf"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type ArchiveController struct {
	pdfs pdf.Repository
}

type archiveControllerFactory interface {
	GetPDFRepository() pdf.Repository
}

func NewArchiveController(factory archiveControllerFactory) *ArchiveController {
	return &ArchiveController{pdfs: factory.GetPDFRepository()}
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
