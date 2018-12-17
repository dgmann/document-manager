package http

import (
	"github.com/dgmann/document-manager/api/repositories"
	"github.com/dgmann/document-manager/api/services"
)

type Factory interface {
	GetRecordRepository() repositories.RecordRepository
	GetImageRepository() repositories.ImageRepository
	GetTagRepository() repositories.TagRepository
	GetPatientRepository() repositories.PatientRepository
	GetCategoryRepository() repositories.CategoryRepository
	GetPDFRepository() repositories.PDFRepository
	GetEventService() *services.EventService
	GetResponseService() *ResponseService
	GetPdfProcessor() (*services.PdfProcessor, error)
}
