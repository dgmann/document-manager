package patients

import (
	"github.com/jmoiron/sqlx"
)

type Patient struct {
	Id   string  `db:"Id"`
	Name *string `db:"Name"`
}

func All(db *sqlx.DB) ([]*Patient, error) {
	var patients []*Patient

	query := `select DISTINCT befund.Pat_Id as Id, max(befund.PatName) as Name
			  from Befund as befund
			  group by befund.Pat_Id
			  ORDER BY befund.Pat_Id`
	err := db.Select(&patients, query)
	if err != nil {
		return nil, err
	}
	return patients, nil
}
