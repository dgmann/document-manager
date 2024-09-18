package m1

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/alexsergivan/transliterator"
	"github.com/dgmann/document-manager/pkg/api"
	go_ora "github.com/sijms/go-ora/v2"
)

type Patient = api.Patient

type Address = api.Address

type DatabaseAdapter struct {
	connectionString string
	db               *sql.DB
	trans            *transliterator.Transliterator
}

func NewDatabaseAdapter(server string, port int, sid string, user string, password string) *DatabaseAdapter {
	connectionString := go_ora.BuildUrl(server, port, "", user, password, map[string]string{
		"SID": sid,
	})
	return NewDatabaseAdapterFromDSN((connectionString))
}

func NewDatabaseAdapterFromDSN(connectionString string) *DatabaseAdapter {
	return &DatabaseAdapter{connectionString: connectionString, trans: transliterator.NewTransliterator(nil)}
}

func (a *DatabaseAdapter) Connect() error {
	db, err := sql.Open("oracle", a.connectionString)
	if err != nil {
		return err
	}
	a.db = db
	return nil
}

func (a *DatabaseAdapter) Close() error {
	return a.db.Close()
}

const query = `
Select Distinct
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
%s
ORDER BY pat.NAME, pat.VORNAME
`

func (a *DatabaseAdapter) GetAllPatients() ([]*Patient, error) {
	rows, err := a.db.Query(fmt.Sprintf(query, "WHERE pat.PATID_EXT IS NOT NULL"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToPatient(rows)
}

func (a *DatabaseAdapter) FindPatientsByName(firstname, lastname string) ([]*Patient, error) {
	name := fmt.Sprintf("%s%%,%s%%", lastname, firstname)
	name = a.trans.Transliterate(name, "de")
	rows, err := a.db.Query(fmt.Sprintf(query, "WHERE pat.PATSNAME like UPPER(:name)"), name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToPatient(rows)
}

func (a *DatabaseAdapter) FindPatientsFuzzy(firstname, lastname string, similarity int, birthDate *time.Time) ([]*Patient, error) {
	name := fmt.Sprintf("%s,%s", lastname, firstname)
	name = a.trans.Transliterate(name, "de")

	whereClause := "WHERE UTL_MATCH.JARO_WINKLER_SIMILARITY(pat.PATSNAME, UPPER(:name)) >= :similarity"
	params := []any{
		name, // required for the select statement part
		name,
		similarity,
	}
	if birthDate != nil {
		whereClause = whereClause + " AND pat.GebDatum = :birthDate"
		params = append(params, birthDate)
	}

	fuzzyQuery := fmt.Sprintf(`
	Select Distinct
	pat.PATID_EXT as PatID,
	pat.Name as Nachname,
	pat.Vorname as Vorname,
	pat.GebDatum as GebDatum,
	wohn.WOHN_PLZ as PLZ,
	wohn.WOHN_STR as Strasse,
	wohnort.ORT_NAME as Ort,
	UTL_MATCH.JARO_WINKLER_SIMILARITY(pat.PATSNAME, UPPER(:name)) as similarity
	From M1PATNT pat
	JOIN M1ADRSS adress on pat.Entty_id = adress.entty_ID
	JOIN M1TELNR telnr on adress.ADRSS_ID = telnr.ADRSS_ID
	LEFT OUTER JOIN M1WOHN wohn on adress.Wohn_ID = wohn.Wohn_ID
	LEFT OUTER JOIN M1ORT wohnort on wohn.ORT_ID = wohnort.ORT_ID
	%s
	ORDER BY similarity DESC, pat.NAME, pat.VORNAME
	`, whereClause)

	rows, err := a.db.Query(fuzzyQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToPatient(rows)
}

func (a *DatabaseAdapter) GetPatient(id string) (*Patient, error) {
	rows, err := a.db.Query(fmt.Sprintf(query, "WHERE pat.PATID_EXT = :id"), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	return rowToPatient(rows)
}

func rowsToPatient(rows *sql.Rows) ([]*Patient, error) {
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

type Scanable interface {
	Scan(dest ...interface{}) error
}

func rowToPatient(row *sql.Rows) (*Patient, error) {
	patient := Patient{}
	address := Address{}
	similarity := 0
	v := []any{
		&patient.Id,
		&patient.LastName,
		&patient.FirstName,
		&patient.BirthDate,
		&address.ZipCode,
		&address.Street,
		&address.City,
	}
	col, err := row.Columns()
	if err != nil {
		return nil, err
	}
	if len(col) == 8 {
		v = append(v, &similarity)
	}
	if err := row.Scan(v...); err != nil {
		return nil, err
	}
	patient.Address = address
	patient.FirstName = strings.TrimSpace(patient.FirstName)
	patient.LastName = strings.TrimSpace(patient.LastName)
	return &patient, err
}
