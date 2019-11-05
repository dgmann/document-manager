package storage

import (
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
	Locate(locatable Locatable) string
	Copy(fromId string, toId string) error
	ModTimeReader
}

type ModTimeReader interface {
	ModTime(resource KeyedResource) (time.Time, error)
}
