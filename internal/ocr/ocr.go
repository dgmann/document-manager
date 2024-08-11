package ocr

import "github.com/dgmann/document-manager/pkg/api"

type PageWithContent struct {
	Id    string
	Image []byte
}

type Client interface {
	// CheckOrientation uses OSD to detect if any page must be rotated.
	// If that is the case, a slice with updated pages is returned.
	// Otherwise, nil is returned indicating that nothing needs to be done.
	CheckOrientation(pages []PageWithContent) ([]api.PageUpdate, error)
	// ExtractText uses OCR to extract text from the provided images.
	// If at least one page with an image is provided, a slice with [api.PageUpdate] is returned.
	// Otherwise, nil is returned indicating that nothing needs to be done.
	ExtractText(pages []PageWithContent) ([]api.PageUpdate, error)
	// Close releases all used resources.
	Close() error
}
