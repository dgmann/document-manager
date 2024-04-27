package storage

import "context"

type Image struct {
	Image  []byte `json:"image"`
	Format string `json:"format"`
}

func NewImage(img []byte, imageType string) *Image {
	return &Image{Image: img, Format: imageType}
}

type ImageService interface {
	ResourceWriter
	ResourceLocator
	GetByRecordId(ctx context.Context, id string) (map[string]*Image, error)
	Copy(ctx context.Context, fromId string, toId string) error
}
