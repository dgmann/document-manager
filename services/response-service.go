package services

import (
	"encoding/json"
	"fmt"
	"github.com/dgmann/document-manager-api/models"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func ToJSON(c *gin.Context, data interface{}) error {
	url := location.Get(c)
	switch v := data.(type) {
	case *models.Record:
		data.(*models.Record).SetURL(url)
	case []*models.Record:
		for _, m := range data.([]*models.Record) {
			m.SetURL(url)
		}
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}
	return json.NewEncoder(c.Writer).Encode(data)
}
