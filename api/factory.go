package main

import (
	"github.com/dgmann/document-manager/api/http"
	"github.com/dgmann/document-manager/api/repositories"
	"github.com/dgmann/document-manager/api/repositories/image"
	"github.com/dgmann/document-manager/api/repositories/pdf"
	"github.com/dgmann/document-manager/api/services"
)

type Factory struct {
	config           *Config
	pdfProcessorUrl  string
	eventService     *services.EventService
	recordRepository repositories.RecordRepository
}

func (f *Factory) GetPdfProcessor() (*services.PdfProcessor, error) {
	return services.NewPDFProcessor(f.pdfProcessorUrl)
}

func (f *Factory) GetResponseService() *http.ResponseService {
	return http.NewResponseService(f.GetImageRepository())
}

func (f *Factory) GetEventService() *services.EventService {
	return f.eventService
}

func (f *Factory) GetRecordRepository() repositories.RecordRepository {
	r := f.config.Db.C("records")
	return repositories.NewDBRecordRepository(r, f.GetImageRepository(), f.GetPDFRepository(), f.GetEventService())
}

func (f *Factory) GetImageRepository() image.Repository {
	return image.NewFileSystemImageRepository(f.config.GetRecordDirectory())
}

func (f *Factory) GetPDFRepository() pdf.Repository {
	return pdf.NewFileSystemPDFRepository(f.config.GetPDFDirectory())
}

func (f *Factory) GetTagRepository() repositories.TagRepository {
	r := f.config.Db.C("records")
	return repositories.NewDBTagRepository(r)
}

func (f *Factory) GetPatientRepository() repositories.PatientRepository {
	patients := f.config.Db.C("patients")
	return repositories.NewDBPatientRepository(patients)
}

func (f *Factory) GetCategoryRepository() repositories.CategoryRepository {
	categories := f.config.Db.C("categories")
	r := f.config.Db.C("records")
	return repositories.NewDBCategoryRepository(categories, r)
}

func NewFactory(config *Config) *Factory {
	fileInfoService := image.NewFileSystemImageRepository(config.GetRecordDirectory())
	eventService := services.NewEventService(fileInfoService)
	eventService.Log()
	f := &Factory{
		config:          config,
		pdfProcessorUrl: config.GetPdfProcessorUrl(),
		eventService:    eventService,
	}
	return f
}
