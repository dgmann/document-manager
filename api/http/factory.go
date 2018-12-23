package http

import (
	"github.com/dgmann/document-manager/api/http/response"
	"github.com/dgmann/document-manager/api/repositories/category"
	"github.com/dgmann/document-manager/api/repositories/image"
	"github.com/dgmann/document-manager/api/repositories/pdf"
	"github.com/dgmann/document-manager/api/repositories/record"
	"github.com/dgmann/document-manager/api/repositories/tag"
	"github.com/dgmann/document-manager/api/services"
)

type Factory interface {
	GetRecordRepository() *record.DatabaseRepository
	GetImageRepository() *image.FileSystemRepository
	GetTagRepository() *tag.DatabaseRepository
	GetCategoryRepository() *category.DatabaseRepository
	GetPDFRepository() *pdf.FileSystemRepository
	GetEventService() *services.EventService
	GetResponseService() *response.Factory
	GetPdfProcessor() (*services.PdfProcessor, error)
}
