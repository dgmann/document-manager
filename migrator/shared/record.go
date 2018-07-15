package shared

import (
	"time"
	"fmt"
)

type PatientRelated interface {
	RecordContainer
	GetPatientId() int
}

type SpezializationRelated interface {
	GetSpezialization() string
}

type Categorizable interface {
	PatientRelated
	SpezializationRelated
}

type RecordContainer interface {
	GetRecord() *Record
}

type Record struct {
	Id             int    `db:"Id"`
	Name           string `db:"Name"`
	PatId          int    `db:"Pat_Id"`
	Spezialization string `db:"Category"`
	Path           string
	SubRecords     []*SubRecord
}

func (r *Record) GetRecord() *Record {
	return r
}

func (r *Record) GetPatientId() int {
	return r.PatId
}

func (r *Record) GetSpezialization() string {
	return r.Spezialization
}

func (r *Record) String() string {
	return fmt.Sprintf("%d, %s: %s", r.PatId, r.Spezialization, r.Path)
}

func (r *Record) Equals(record *Record) bool {
	return r.PatId == record.PatId &&
		r.Spezialization == record.Spezialization &&
		r.Name == record.Name
}

type SubRecord struct {
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
	Path          string
}
