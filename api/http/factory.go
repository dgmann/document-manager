package http

import (
	"github.com/dgmann/document-manager/api/repositories/category"
	"github.com/dgmann/document-manager/api/repositories/image"
	"github.com/dgmann/document-manager/api/repositories/patient"
	"github.com/dgmann/document-manager/api/repositories/pdf"
	"github.com/dgmann/document-manager/api/repositories/record"
	"github.com/dgmann/document-manager/api/repositories/tag"
	"github.com/dgmann/document-manager/api/services"
)

type Factory interface {
	GetRecordRepository() record.Repository
	GetImageRepository() image.Repository
	GetTagRepository() tag.Repository
	GetPatientRepository() patient.Repository
	GetCategoryRepository() category.Repository
	GetPDFRepository() pdf.Repository
	GetEventService() *services.EventService
	GetResponseService() *ResponseService
	GetPdfProcessor() (*services.PdfProcessor, error)
}
