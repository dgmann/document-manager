package pdf

import (
	"github.com/dgmann/document-manager/pdf-processor/pdfprocessor"
	"bytes"
	_ "image/jpeg"
	_ "image/gif"
	_ "image/png"
)

type Result []pdfprocessor.ImageResult

func (r Result) ToImages() []*bytes.Buffer {
	images := make([]*bytes.Buffer, len(r))
	for _, element := range r {
		images[element.PageNumber] = bytes.NewBuffer(element.Image)
	}
	return images
}
