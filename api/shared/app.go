package shared

import (
	"github.com/dgmann/document-manager/api/pdf"
	"github.com/dgmann/document-manager/api/repositories"
)

type App struct {
	Records      repositories.RecordRepository
	Images       repositories.ImageRepository
	Tags         *repositories.TagRepository
	Patients     *repositories.PatientRepository
	Categories   *repositories.CategoryRepository
	PDFProcessor *pdf.PDFProcessor
}
