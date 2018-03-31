package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/models"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func RespondAsJSON(c *gin.Context, data interface{}) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	url := location.Get(c)
	switch data.(type) {
	case *models.Record:
		data.(*models.Record).SetURL(url)
	case []*models.Record:
		for _, m := range data.([]*models.Record) {
			m.SetURL(url)
		}
	}
	if err := json.NewEncoder(c.Writer).Encode(data); err != nil {
		c.AbortWithError(400, err)
	}
}
