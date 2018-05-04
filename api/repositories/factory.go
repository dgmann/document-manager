package repositories

type factory struct {
	config RepositoryConfig
}

type Factory interface {
	GetRecordRepository() RecordRepository
	GetImageRepository() ImageRepository
	GetTagRepository() TagRepository
	GetPatientRepository() PatientRepository
	GetCategoryRepository() CategoryRepository
}

func NewFactory(config RepositoryConfig) Factory {
	return &factory{config: config}
}

func (f *factory) GetRecordRepository() RecordRepository {
	return NewRecordRepository(f.config)
}

func (f *factory) GetImageRepository() ImageRepository {
	return NewImageReporitory(f.config)
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
