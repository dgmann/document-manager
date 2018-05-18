package m1

import (
	_ "gopkg.in/rana/ora.v4"
	"database/sql"
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
	db, err := sql.Open("ora", a.connectionString)
	if err != nil {
		return err
	}
	a.db = db
	return nil
}

func (a *DatabaseAdapter) Close() error {
	return a.db.Close()
}

func (a *DatabaseAdapter) GetPatient(id string) (*Patient, error) {
	patient := Patient{}

	err := a.db.QueryRow(`Select Distinct 
									pat.PATID_EXT as PatID, 
									pat.Name as Nachname, 
									pat.Vorname as Vorname,
									pat.GebDatum as GebDatum,
									telort.ORT_VORWAHL as Vorwahl,
									ansch.ANSCH_NR as Nummer,
									WOHN.WOHN_PLZ as PLZ,
									WOHN.WOHN_STR as Strasse,
									wohnort.ORT_NAME as Ort,
									CONCAT(telort.ORT_VORWAHL,ansch.ANSCH_NR) as FULLNumber
									From M1PATNT pat 
									JOIN M1ADRSS adress on pat.Entty_id = adress.entty_ID 
									JOIN M1TELNR telnr on adress.ADRSS_ID = telnr.ADRSS_ID 
									JOIN M1ANSCH ansch on telnr.ansch_id = ansch.ansch_id 
									LEFT OUTER JOIN M1ORT telort on ansch.ORT_ID = telort.ORT_ID
									LEFT OUTER JOIN M1WOHN wohn on adress.Wohn_ID = wohn.Wohn_ID
									LEFT OUTER JOIN M1ORT wohnort on wohn.ORT_ID = wohnort.ORT_ID
									WHERE pat.PATID_EXT = :id 
									ORDER BY pat.NAME, pat.VORNAME`, id).Scan(&patient.Id, &patient.LastName, &patient.FirstName, &patient.BirthDate)
	if err != nil {
		return nil, err
	}

	return &patient, nil
}
