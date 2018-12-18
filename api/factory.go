package main

import (
	"github.com/dgmann/document-manager/api/http/response"
	"github.com/dgmann/document-manager/api/repositories/category"
	"github.com/dgmann/document-manager/api/repositories/image"
	"github.com/dgmann/document-manager/api/repositories/patient"
	"github.com/dgmann/document-manager/api/repositories/pdf"
	"github.com/dgmann/document-manager/api/repositories/record"
	"github.com/dgmann/document-manager/api/repositories/tag"
	"github.com/dgmann/document-manager/api/services"
)

type Factory struct {
	config           *Config
	pdfProcessorUrl  string
	eventService     *services.EventService
	recordRepository record.Repository
}

func (f *Factory) GetPdfProcessor() (*services.PdfProcessor, error) {
	return services.NewPDFProcessor(f.pdfProcessorUrl)
}

func (f *Factory) GetResponseService() *response.Factory {
	return response.NewFactory(f.GetImageRepository())
}

func (f *Factory) GetEventService() *services.EventService {
	return f.eventService
}

func (f *Factory) GetRecordRepository() *record.DatabaseRepository {
	r := f.config.Db.C("records")
	return record.NewDatabaseRepository(r, f.GetImageRepository(), f.GetPDFRepository(), f.GetEventService())
}

func (f *Factory) GetImageRepository() *image.FileSystemRepository {
	return image.NewFileSystemRepository(f.config.GetRecordDirectory())
}

func (f *Factory) GetPDFRepository() *pdf.FileSystemRepository {
	return pdf.NewFileSystemRepository(f.config.GetPDFDirectory())
}

func (f *Factory) GetTagRepository() *tag.DatabaseRepository {
	r := f.config.Db.C("records")
	return tag.NewDatabaseRepository(r)
}

func (f *Factory) GetPatientRepository() *patient.DatabaseRepository {
	patients := f.config.Db.C("patients")
	return patient.NewDatabaseRepository(patients)
}

func (f *Factory) GetCategoryRepository() *category.DatabaseRepository {
	categories := f.config.Db.C("categories")
	r := f.config.Db.C("records")
	return category.NewDatabaseRepository(categories, r)
}

func NewFactory(config *Config) *Factory {
	fileInfoService := image.NewFileSystemRepository(config.GetRecordDirectory())
	eventService := services.NewEventService(fileInfoService)
	eventService.Log()
	f := &Factory{
		config:          config,
		pdfProcessorUrl: config.GetPdfProcessorUrl(),
		eventService:    eventService,
	}
	return f
}
