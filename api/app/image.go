package app

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Image struct {
	Image  []byte `json:"image"`
	Format string `json:"format"`
}

func NewImage(img []byte, imageType string) *Image {
	return &Image{Image: img, Format: imageType}
}

type ImageService interface {
	ResourceWriter
	Get(id string) (map[string]*Image, error)
	Serve(context *gin.Context, recordId string, imageId string, format string)
	Copy(fromId string, toId string) error
	ModTimeReader
}

type ModTimeReader interface {
	ModTime(resource KeyedResource) (time.Time, error)
}
