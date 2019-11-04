package m1

import (
	"database/sql"
	_ "github.com/mattn/go-oci8"
	"strings"
)

type Adapter interface {
}

type DatabaseAdapter struct {
	connectionString string
	db               *sql.DB
}

func NewDatabaseAdapter(connectionString string) *DatabaseAdapter {
	return &DatabaseAdapter{connectionString: connectionString}
}

func (a *DatabaseAdapter) Connect() error {
	db, err := sql.Open("oci8", a.connectionString)
	if err != nil {
		return err
	}
	a.db = db
	return nil
}

func (a *DatabaseAdapter) Close() error {
	return a.db.Close()
}

func (a *DatabaseAdapter) GetAllPatients() ([]*Patient, error) {
	rows, err := a.db.Query(`Select Distinct
								pat.PATID_EXT as PatID,
								pat.Name as Nachname,
								pat.Vorname as Vorname,
								pat.GebDatum as GebDatum,
								wohn.WOHN_PLZ as PLZ,
								wohn.WOHN_STR as Strasse,
								wohnort.ORT_NAME as Ort
								From M1PATNT pat
								JOIN M1ADRSS adress on pat.Entty_id = adress.entty_ID
								JOIN M1TELNR telnr on adress.ADRSS_ID = telnr.ADRSS_ID
								LEFT OUTER JOIN M1WOHN wohn on adress.Wohn_ID = wohn.Wohn_ID
								LEFT OUTER JOIN M1ORT wohnort on wohn.ORT_ID = wohnort.ORT_ID
								ORDER BY pat.NAME, pat.VORNAME`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients = make([]*Patient, 0)
	for rows.Next() {
		patient, err := rowToPatient(rows)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}
	return patients, nil
}

func (a *DatabaseAdapter) FindPatientsByName(firstname, lastname string) ([]*Patient, error) {
	firstname += "%"
	lastname += "%"
	rows, err := a.db.Query(`Select Distinct
								pat.PATID_EXT as PatID,
								pat.Name as Nachname,
								pat.Vorname as Vorname,
								pat.GebDatum as GebDatum,
								wohn.WOHN_PLZ as PLZ,
								wohn.WOHN_STR as Strasse,
								wohnort.ORT_NAME as Ort
								From M1PATNT pat
								JOIN M1ADRSS adress on pat.Entty_id = adress.entty_ID
								JOIN M1TELNR telnr on adress.ADRSS_ID = telnr.ADRSS_ID
								LEFT OUTER JOIN M1WOHN wohn on adress.Wohn_ID = wohn.Wohn_ID
								LEFT OUTER JOIN M1ORT wohnort on wohn.ORT_ID = wohnort.ORT_ID
								WHERE LOWER(pat.Vorname) LIKE LOWER(:firstname)
								AND LOWER(pat.Name) LIKE LOWER(:lastname)
								ORDER BY pat.NAME, pat.VORNAME`, firstname, lastname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients = make([]*Patient, 0)
	for rows.Next() {
		patient, err := rowToPatient(rows)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}
	return patients, nil
}

func (a *DatabaseAdapter) GetPatient(id string) (*Patient, error) {
	row := a.db.QueryRow(`Select Distinct
								pat.PATID_EXT as PatID,
								pat.Name as Nachname,
								pat.Vorname as Vorname,
								pat.GebDatum as GebDatum,
								wohn.WOHN_PLZ as PLZ,
								wohn.WOHN_STR as Strasse,
								wohnort.ORT_NAME as Ort
								From M1PATNT pat
								JOIN M1ADRSS adress on pat.Entty_id = adress.entty_ID
								JOIN M1TELNR telnr on adress.ADRSS_ID = telnr.ADRSS_ID
								LEFT OUTER JOIN M1WOHN wohn on adress.Wohn_ID = wohn.Wohn_ID
								LEFT OUTER JOIN M1ORT wohnort on wohn.ORT_ID = wohnort.ORT_ID
								WHERE pat.PATID_EXT = :id
								ORDER BY pat.NAME, pat.VORNAME`, id)
	return rowToPatient(row)
}

type Scanable interface {
	Scan(dest ...interface{}) error
}

func rowToPatient(row Scanable) (*Patient, error) {
	patient := Patient{}
	address := Address{}
	err := row.Scan(
		&patient.Id,
		&patient.LastName,
		&patient.FirstName,
		&patient.BirthDate,
		&address.ZipCode,
		&address.Street,
		&address.City)
	patient.Address = address
	patient.FirstName = strings.TrimSpace(patient.FirstName)
	patient.LastName = strings.TrimSpace(patient.LastName)
	return &patient, err
}
