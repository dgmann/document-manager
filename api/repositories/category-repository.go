package repositories

type CategoryRepository struct {
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (t *CategoryRepository) All() ([]string, error) {
	return []string{
		"Innere Medizin",
		"Chirugie",
		"Neurologie",
	}, nil
}
