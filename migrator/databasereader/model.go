package databasereader

import "time"

type PatientCategoryRecordMap map[int]RecordsByCategory

type RecordsByCategory map[string]*Record

type Record struct {
	Id       int    `db:"Id"`
	PatId    int    `db:"Pat_Id"`
	Category string `db:"Category"`
	PdfFiles []*PdfFile
}

type PdfFile struct {
	Id         int        `db:"Id"`
	Date       *time.Time `db:"Date"`
	ReceivedAt time.Time  `db:"Timestamp"`
	State      int        `db:"State"`
	Type       int        `db:"Type"`
	SenderNr   *string    `db:"SenderNr"`
	Path       string
	BefundId   *int       `db:"Befund_Id"`
}
