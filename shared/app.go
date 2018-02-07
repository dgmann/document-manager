package shared

import (
	"github.com/dgmann/document-manager-api/repositories"
	"github.com/dgmann/document-manager-api/services"
)

type App struct {
	Records *repositories.RecordRepository
	Images  *repositories.FileSystemImageRepository
	PDFProcessor *services.PDFProcessor
}
