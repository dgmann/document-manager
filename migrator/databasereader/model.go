package databasereader

import (
	"time"
	"path"
)

type DatabaseIndex struct {
	data              PatientCategoryRecordMap
	TotalRecordCount  int
	TotalPatientCount int
}

func NewDatabaseIndex(data PatientCategoryRecordMap) *DatabaseIndex {
	totalRecordCount := 0
	for _, recordsByCategory := range data {
		//Because each spezialization has only one record, we can count only the spezializations
		totalRecordCount += len(recordsByCategory)
	}
	return &DatabaseIndex{data: data, TotalRecordCount: totalRecordCount, TotalPatientCount: len(data)}
}

type PatientCategoryRecordMap map[int]RecordsByCategory

type RecordsByCategory map[string]*Record

type Record struct {
	Id       int    `db:"Id"`
	PatId    int    `db:"Pat_Id"`
	Category string `db:"Category"`
	PdfFiles []*PdfFile
}

type PdfFile struct {
	Id            int        `db:"Id"`
	Name          string     `db:"Name"`
	Date          *time.Time `db:"Date"`
	ReceivedAt    time.Time  `db:"Timestamp"`
	State         int        `db:"State"`
	Type          *string    `db:"Type"`
	SenderNr      *string    `db:"SenderNr"`
	PathExtension *string    `db:"PathExtension"`
	BefundId      *int       `db:"Befund_Id"`
	PatId         *int       `db:"Pat_Id"`
}

func (p *PdfFile) GetPath() string {
	return path.Join(*p.PathExtension, p.Name)
}
