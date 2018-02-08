package shared

import (
	"github.com/dgmann/document-manager-api/repositories"
	"github.com/dgmann/document-manager-api/pdf"
)

type App struct {
	Records *repositories.RecordRepository
	Images  *repositories.FileSystemImageRepository
	Tags *repositories.TagRepository
	PDFProcessor *pdf.PDFProcessor
}
