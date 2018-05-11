package repositories

import "github.com/dgmann/document-manager/api/services"

type factory struct {
	config       RepositoryConfig
	eventService *services.EventService
}

type Factory interface {
	GetRecordRepository() RecordRepository
	GetImageRepository() ImageRepository
	GetTagRepository() TagRepository
	GetPatientRepository() PatientRepository
	GetCategoryRepository() CategoryRepository
	GetPDFRepository() PDFRepository
}

func NewFactory(config RepositoryConfig, eventService *services.EventService) Factory {
	return &factory{config: config, eventService: eventService}
}

func (f *factory) GetRecordRepository() RecordRepository {
	return NewRecordRepository(f.config, f.eventService)
}

func (f *factory) GetImageRepository() ImageRepository {
	return NewImageRepository(f.config)
}

func (f *factory) GetTagRepository() TagRepository {
	return NewTagRepository(f.config)
}

func (f *factory) GetPatientRepository() PatientRepository {
	return NewPatientRepository(f.config)
}

func (f *factory) GetCategoryRepository() CategoryRepository {
	return NewCategoryRepository(f.config)
}

func (f *factory) GetPDFRepository() PDFRepository {
	return NewPDFRepository(f.config)
}
