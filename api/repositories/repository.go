package repositories

import (
	"github.com/dgmann/document-manager/api/shared"
	"github.com/dgmann/document-manager/api/services"
)

type RepositoryConfig interface {
	shared.DatabaseConfig
	shared.RecordDirectoryConfig
	shared.PDFDirectoryConfig
}

func NewRecordRepository(config RepositoryConfig, eventService *services.EventService) RecordRepository {
	records := config.GetDatabase().C("records")
	images := NewImageRepository(config)
	pdfs := NewPDFRepository(config)
	return newDBRecordRepository(records, images, pdfs, eventService)
}

func NewImageRepository(config shared.RecordDirectoryConfig) ImageRepository {
	return newFileSystemImageRepository(config.GetRecordDirectory())
}

func NewTagRepository(config shared.DatabaseConfig) TagRepository {
	records := config.GetDatabase().C("records")
	return newDBTagRepository(records)
}

func NewPatientRepository(config shared.DatabaseConfig) PatientRepository {
	patients := config.GetDatabase().C("patients")
	return newDBPatientRepository(patients)
}

func NewCategoryRepository(config shared.DatabaseConfig) CategoryRepository {
	categories := config.GetDatabase().C("categories")
	records := config.GetDatabase().C("records")
	return newDBCategoryRepository(categories, records)
}

func NewPDFRepository(config shared.PDFDirectoryConfig) PDFRepository {
	return newFileSystemPDFRepository(config.GetPDFDirectory())
}
