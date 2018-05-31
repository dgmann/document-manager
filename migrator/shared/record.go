package shared

import (
	"time"
	"fmt"
)

type PatientRelated interface {
	GetPatientId() int
}

type SpezializationRelated interface {
	GetSpezialization() string
}

type Categorizable interface {
	PatientRelated
	SpezializationRelated
}

type Record struct {
	Id             int    `db:"Id"`
	Name           string
	PatId          int    `db:"Pat_Id"`
	Spezialization string `db:"Category"`
	Path           string
	SubRecords     []*SubRecord
}

func (r *Record) GetPatientId() int {
	return r.PatId
}

func (r *Record) GetSpezialization() string {
	return r.Spezialization
}

func (r *Record) String() string {
	return fmt.Sprintf("%s, %s: %s", r.PatId, r.Spezialization, r.Path)
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
