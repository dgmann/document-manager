package categories

import "github.com/jmoiron/sqlx"

type Category struct {
	Id   string `db:"Id"`
	Name string `db:"Name"`
}

func All(db *sqlx.DB) ([]*Category, error) {
	var categories []*Category

	query := `select Name as Id, Description as Name 
			  from Spezialisations`
	err := db.Select(&categories, query)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
