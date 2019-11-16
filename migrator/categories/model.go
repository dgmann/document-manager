package categories

import (
	"fmt"
	"github.com/dgmann/document-manager/api/datastore"
	"github.com/jmoiron/sqlx"
)

func All(db *sqlx.DB) ([]datastore.Category, error) {
	var categories []datastore.Category

	query := `select Name, Description
			  from Spezialisations`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error execution categories query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category datastore.Category
		if err := rows.Scan(&category.Id, &category.Name); err != nil {
			return nil, fmt.Errorf("error fetching database row: %w", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}
