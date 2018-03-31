package pdf

import (
	"github.com/dgmann/document-manager/shared"
	"bytes"
	_ "image/jpeg"
	_ "image/gif"
	_ "image/png"
)

type Result []shared.ImageResult

func (r Result) ToImages() []*bytes.Buffer {
	images := make([]*bytes.Buffer, len(r))
	for _, element := range r {
		images[element.PageNumber] = bytes.NewBuffer(element.Image)
	}
	return images
}
